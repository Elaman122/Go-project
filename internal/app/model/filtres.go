package model

import (
	"math"
	"strings"

	"github.com/Elaman122/Go-project/internal/app/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
	Offset   	 int
}

// Metadata хранит метаданные пагинации.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "должно быть больше 0")
	v.Check(f.Page <= 10_000_0000, "page", "должно быть максимум 10 миллионов")
	v.Check(f.PageSize > 0, "page_size", "должно быть больше 0")
	v.Check(f.PageSize <= 100, "page_size", "должно быть максимум 100")

	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

func (f Filters) sortColumn() string {
    for _, safeValue := range f.SortSafeList {
        if f.Sort == safeValue || f.Sort == "-"+safeValue {
            return strings.TrimPrefix(f.Sort, "-")
        }
    }

    // Вернуть значение по умолчанию, например, первое значение из списка безопасных значений
    return f.SortSafeList[0]
}



func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// NewFilters создает новый объект Filters с значениями по умолчанию.
func NewFilters() Filters {
    return Filters{
        Page:         1,
        PageSize:     10,
        Sort:         "",
        SortSafeList: []string{"rate", "code", "timestamp"}, // Пример безопасных значений для сортировки
    }
}
