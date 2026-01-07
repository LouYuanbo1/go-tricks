package repo

import (
	"context"
	"fmt"
	"go-tricks/gorm/decoupling-method/model"
	"log"

	"gorm.io/gorm"
)

type genericRepo[T any, PT model.PointerModel[T]] struct {
	db *gorm.DB
}

func NewGenericRepo[T any, PT model.PointerModel[T]](db *gorm.DB) GenericRepo[T, PT] {
	return &genericRepo[T, PT]{db: db}
}

func (r *genericRepo[T, PT]) Create(ctx context.Context, ptrModel PT) error {
	if ptrModel == nil {
		var model T
		ptr := PT(&model)
		return fmt.Errorf("create %s failed, ptrModel is nil", ptr.TableName())
	}

	fmt.Printf("type ptrModel: %T", ptrModel)

	result := r.db.WithContext(ctx).
		Create(ptrModel)
	if result.Error != nil {
		log.Printf("create %s failed, error: %v", ptrModel.TableName(), result.Error)
		return fmt.Errorf("create %s failed, error: %v", ptrModel.TableName(), result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("create %s failed, no rows affected", ptrModel.TableName())
		return fmt.Errorf("create %s failed, no rows affected", ptrModel.TableName())
	}
	return nil
}

func (r *genericRepo[T, PT]) CreateInBatches(ctx context.Context, ptrModels []PT, batchSize int) error {
	// 检查ptrModels是否为空
	if len(ptrModels) == 0 || batchSize <= 0 {
		var model T
		ptr := PT(&model)
		log.Printf("create %s in batchs failed, no models provided", ptr.TableName())
		return nil
	}
	result := r.db.WithContext(ctx).
		CreateInBatches(ptrModels, batchSize)
	if result.Error != nil {
		log.Printf("create %s in batchs failed, error: %v", ptrModels[0].TableName(), result.Error)
		return fmt.Errorf("create %s in batchs failed, error: %v", ptrModels[0].TableName(), result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("create %s in batchs failed, no rows affected", ptrModels[0].TableName())
		return fmt.Errorf("create %s in batchs failed, no rows affected", ptrModels[0].TableName())
	}
	return nil
}

func (r *genericRepo[T, PT]) GetByID(ctx context.Context, id uint64) (PT, error) {
	var model T
	ptrModel := PT(&model)

	if id == 0 {
		log.Printf("get %s by id %d failed, id must be greater than 0", ptrModel.TableName(), id)
		return nil, nil
	}

	result := r.db.WithContext(ctx).
		Where(fmt.Sprintf("%s = ?", ptrModel.GetPrimaryKey()), id).
		First(ptrModel)
	if result.Error != nil {
		log.Printf("get %s by id %d failed, error: %v", ptrModel.TableName(), id, result.Error)
		return nil, fmt.Errorf("get %s by id %d failed, error: %v", ptrModel.TableName(), id, result.Error)
	}
	return ptrModel, nil
}

func (r *genericRepo[T, PT]) GetByIDs(ctx context.Context, ids []uint64) ([]PT, error) {
	var model T
	ptrModel := PT(&model)

	if len(ids) == 0 {
		log.Printf("get %s by ids %v failed, no ids provided", ptrModel.TableName(), ids)
		return nil, nil
	}

	fmt.Printf("type model: %T, type ptrModel: %T", model, ptrModel)

	ptrModels := make([]PT, 0, len(ids))

	fmt.Printf("type ptrModels: %T", ptrModels)

	result := r.db.WithContext(ctx).
		Where(fmt.Sprintf("%s IN ?", ptrModel.GetPrimaryKey()), ids).
		Find(&ptrModels)
	if result.Error != nil {
		log.Printf("get %s by ids %v failed, error: %v", ptrModel.TableName(), ids, result.Error)
		return nil, fmt.Errorf("get %s by ids %v failed, error: %v", ptrModel.TableName(), ids, result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("get %s by ids %v failed, no rows affected", ptrModel.TableName(), ids)
		return nil, fmt.Errorf("get %s by ids %v failed, no rows affected", ptrModel.TableName(), ids)
	}
	return ptrModels, nil
}

func (r *genericRepo[T, PT]) GetByPage(ctx context.Context, page, pageSize int) ([]PT, error) {
	var model T
	ptrModel := PT(&model)

	if page <= 0 || pageSize <= 0 {
		log.Printf("get %s by page %d, pageSize %d failed, page and pageSize must be greater than 0", ptrModel.TableName(), page, pageSize)
		return nil, nil
	}

	ptrModels := make([]PT, 0, pageSize)

	result := r.db.WithContext(ctx).
		Order(fmt.Sprintf("%s ASC", ptrModel.GetPrimaryKey())).
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&ptrModels)
	if result.Error != nil {
		log.Printf("get %s by page %d, pageSize %d failed, error: %v", ptrModel.TableName(), page, pageSize, result.Error)
		return nil, fmt.Errorf("get %s by page %d, pageSize %d failed, error: %v", ptrModel.TableName(), page, pageSize, result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("get %s by page %d, pageSize %d failed, no rows affected", ptrModel.TableName(), page, pageSize)
		return nil, fmt.Errorf("get %s by page %d, pageSize %d failed, no rows affected", ptrModel.TableName(), page, pageSize)
	}
	return ptrModels, nil
}

func (r *genericRepo[T, PT]) GetByCursor(ctx context.Context, cursor uint64, pageSize int) ([]PT, uint64, bool, error) {
	var model T
	ptrModel := PT(&model)

	if pageSize <= 0 {
		log.Printf("get %s by cursor %d, pageSize %d failed, pageSize must be greater than 0", ptrModel.TableName(), cursor, pageSize)
		return nil, cursor, false, nil
	}
	limit := pageSize + 1

	ptrModels := make([]PT, 0, pageSize)
	result := r.db.WithContext(ctx).
		Where(fmt.Sprintf("%s = 0 OR %s > ?", ptrModel.GetPrimaryKey(), ptrModel.GetPrimaryKey()), cursor).
		Order(fmt.Sprintf("%s ASC", ptrModel.GetPrimaryKey())).
		Limit(int(limit)).
		Find(&ptrModels)
	if result.Error != nil {
		log.Printf("get %s by cursor %d, pageSize %d failed, error: %v", ptrModel.TableName(), cursor, pageSize, result.Error)
		return nil, cursor, false, fmt.Errorf("get %s by cursor %d, pageSize %d failed, error: %v", ptrModel.TableName(), cursor, pageSize, result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("get %s by cursor %d, pageSize %d failed, no rows affected", ptrModel.TableName(), cursor, pageSize)
		return nil, cursor, false, fmt.Errorf("get %s by cursor %d, pageSize %d failed, no rows affected", ptrModel.TableName(), cursor, pageSize)
	}
	hasMore := len(ptrModels) > pageSize
	if hasMore {
		ptrModels = ptrModels[:pageSize]
	}
	newCursor := cursor
	if len(ptrModels) > 0 {
		newCursor = ptrModels[len(ptrModels)-1].GetID()
	}
	return ptrModels, newCursor, hasMore, nil
}

func (r *genericRepo[T, PT]) Update(ctx context.Context, ptrModel PT) error {
	if ptrModel == nil {
		var model T
		ptr := PT(&model)
		log.Printf("update %s failed, ptrModel is nil", ptr.TableName())
		return nil
	}
	fmt.Printf("type ptrModel: %T", ptrModel)

	result := r.db.WithContext(ctx).
		Updates(ptrModel)
	if result.Error != nil {
		log.Printf("update %s failed, error: %v", ptrModel.TableName(), result.Error)
		return fmt.Errorf("update %s failed, error: %v", ptrModel.TableName(), result.Error)
	}
	return nil
}

func (r *genericRepo[T, PT]) DeleteByID(ctx context.Context, id uint64) error {
	var model T
	pt := PT(&model)

	if id == 0 {
		log.Printf("delete %s failed, id is 0", pt.TableName())
		return nil
	}
	result := r.db.WithContext(ctx).
		Where(fmt.Sprintf("%s = ?", pt.GetPrimaryKey()), id).
		Delete(pt)
	if result.Error != nil {
		log.Printf("delete %s failed, error: %v", pt.TableName(), result.Error)
		return fmt.Errorf("delete %s failed, error: %v", pt.TableName(), result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("delete %s failed, no rows affected", pt.TableName())
		return fmt.Errorf("delete %s failed, no rows affected", pt.TableName())
	}
	return nil
}

func (r *genericRepo[T, PT]) DeleteByIDs(ctx context.Context, ids []uint64) error {
	var model T
	pt := PT(&model)

	if len(ids) == 0 {
		log.Printf("delete %s by ids failed, no ids provided", pt.TableName())
		return nil
	}
	result := r.db.WithContext(ctx).
		Where(fmt.Sprintf("%s IN ?", pt.GetPrimaryKey()), ids).
		Delete(pt)
	if result.Error != nil {
		log.Printf("delete %s by ids %v failed, error: %v", pt.TableName(), ids, result.Error)
		return fmt.Errorf("delete %s by ids %v failed, error: %v", pt.TableName(), ids, result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("delete %s by ids %v failed, no rows affected", pt.TableName(), ids)
		return fmt.Errorf("delete %s by ids %v failed, no rows affected", pt.TableName(), ids)
	}
	return nil
}
