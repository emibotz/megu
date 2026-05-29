package todo

import (
	"context"
	"fmt"
	"time"
)

// 待办服务
type Service struct {
	store Store
}

// 创建待办服务，需传入待办仓库
func NewService(
	store Store,
) *Service {
	return &Service{
		store: store,
	}
}

// 创建待办，需传入待办标题和内容
func (s *Service) Create(ctx context.Context, title string, content string) error {
	td := New(title, content)

	if err := s.store.Create(ctx, td); err != nil {
		return err
	}

	return nil
}

// 获取待办信息，页码从 0 开始
func (s *Service) List(ctx context.Context, pageNum int) ([]ListResp, error) {
	const limit = 5
	offset := pageNum * limit

	return s.store.List(ctx, limit, offset)
}

// 获取指定序号的待办的详细信息
func (s *Service) GetByNumber(ctx context.Context, num int) (*TODO, error) {
	td, err := s.store.GetByNumber(ctx, num)
	if err != nil {
		return nil, err
	}
	if td == nil {
		return nil, fmt.Errorf("no item with number %d", num)
	}
	return td, nil
}

// 设置指定序号的待办的标题
func (s *Service) SetTitle(ctx context.Context, num int, title string) error {
	td, err := s.GetByNumber(ctx, num)
	if err != nil {
		return err
	}

	td.Title = title

	if err := s.store.Update(ctx, td); err != nil {
		return err
	}

	return nil
}

// 设置指定序号的待办的内容
func (s *Service) SetContent(ctx context.Context, num int, content string) error {
	td, err := s.GetByNumber(ctx, num)
	if err != nil {
		return err
	}

	td.Content = content

	if err := s.store.Update(ctx, td); err != nil {
		return err
	}

	return nil
}

// 设置指定序号的待办是否完成
func (s *Service) SetCompleted(ctx context.Context, num int, completed bool) error {
	td, err := s.GetByNumber(ctx, num)
	if err != nil {
		return err
	}

	td.Completed = completed

	if err := s.store.Update(ctx, td); err != nil {
		return err
	}

	return nil
}

// 设置指定序号的待办是否放入回收站
func (s *Service) SetDeleted(ctx context.Context, num int, deleted bool) error {
	td, err := s.GetByNumber(ctx, num)
	if err != nil {
		return err
	}

	if deleted && td.DeletedAt == nil {
		now := time.Now()
		td.DeletedAt = &now
	} else if !deleted {
		td.DeletedAt = nil
	}

	if err := s.store.Update(ctx, td); err != nil {
		return err
	}

	return nil
}

// 永久删除指定序号的待办
func (s *Service) Destroy(ctx context.Context, num int) error {
	td, err := s.GetByNumber(ctx, num)
	if err != nil {
		return err
	}

	if err := s.store.Delete(ctx, td); err != nil {
		return err
	}

	return nil
}
