package todo

import "context"

type ListResp struct {
	TD  TODO
	Num int
}

type Store interface {
	Create(ctx context.Context, td *TODO) error

	List(ctx context.Context, limit int, offset int) ([]ListResp, error)
	// 通过是否完成筛选列出的待办项
	// ListByCompleted(ctx context.Context, completed bool, limit int, offset int) ([]ListResp, error)
	// 获取指定序号的待办，若没有指定记录则返回 (nil, nil)
	GetByNumber(ctx context.Context, num int) (*TODO, error)

	Update(ctx context.Context, td *TODO) error

	Delete(ctx context.Context, td *TODO) error
}
