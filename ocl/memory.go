// +build cl11 cl12

package ocl

import (
	"errors"
	"gocl/cl"
	"unsafe"
)

type Memory interface {
	GetID() cl.CL_mem
	GetInfo(param_name cl.CL_mem_info) (interface{}, error)
	Retain() error
	Release() error

	SetCallback(pfn_notify cl.CL_mem_notify,
		user_data unsafe.Pointer) error
	EnqueueUnmap(queue CommandQueue,
		mapped_ptr unsafe.Pointer,
		event_wait_list []Event) (Event, error)
}

type memory struct {
	memory_id cl.CL_mem
}

func (this *memory) GetID() cl.CL_mem {
	return this.memory_id
}

func (this *memory) GetInfo(param_name cl.CL_mem_info) (interface{}, error) {
	/* param data */
	var param_value interface{}
	var param_size cl.CL_size_t
	var errCode cl.CL_int

	/* Find size of param data */
	if errCode = cl.CLGetMemObjectInfo(this.memory_id, param_name, 0, nil, &param_size); errCode != cl.CL_SUCCESS {
		return nil, errors.New("GetInfo failure with errcode_ret " + string(errCode))
	}

	/* Access param data */
	if errCode = cl.CLGetMemObjectInfo(this.memory_id, param_name, param_size, &param_value, nil); errCode != cl.CL_SUCCESS {
		return nil, errors.New("GetInfo failure with errcode_ret " + string(errCode))
	}

	return param_value, nil
}

func (this *memory) Retain() error {
	if errCode := cl.CLRetainMemObject(this.memory_id); errCode != cl.CL_SUCCESS {
		return errors.New("Retain failure with errcode_ret " + string(errCode))
	}
	return nil
}

func (this *memory) Release() error {
	if errCode := cl.CLReleaseMemObject(this.memory_id); errCode != cl.CL_SUCCESS {
		return errors.New("Release failure with errcode_ret " + string(errCode))
	}
	return nil
}

func (this *memory) SetCallback(pfn_notify cl.CL_mem_notify, user_data unsafe.Pointer) error {
	if errCode := cl.CLSetMemObjectDestructorCallback(this.memory_id, pfn_notify, user_data); errCode != cl.CL_SUCCESS {
		return errors.New("SetCallback failure with errcode_ret " + string(errCode))
	} else {
		return nil
	}
}

func (this *memory) EnqueueUnmap(queue CommandQueue, mapped_ptr unsafe.Pointer, event_wait_list []Event) (Event, error) {
	var event_id cl.CL_event

	numEvents := cl.CL_uint(len(event_wait_list))
	events := make([]cl.CL_event, numEvents)
	for i := cl.CL_uint(0); i < numEvents; i++ {
		events[i] = event_wait_list[i].GetID()
	}

	if errCode := cl.CLEnqueueUnmapMemObject(queue.GetID(), this.memory_id, mapped_ptr, numEvents, events, &event_id); errCode != cl.CL_SUCCESS {
		return nil, errors.New("EnqueueMarkerWithWaitList failure with errcode_ret " + string(errCode))
	} else {
		return &event{event_id}, nil
	}
}
