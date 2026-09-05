package helper

import (
	"strconv"
	"strings"

	"gakuren-system.com/pkg/db"
	"gakuren-system.com/pkg/model"
)

func SearchGender(tenantUUID string, payload model.SearchPayload) ([]model.ReadGenderResult, *model.DataStatistics, error) {
	var param []interface{}

	//* base query
	query := `
	with datas as(
		select
			g."uuid" 
			,g."name" 
			,g.abbr_name 
			,s."name" as status
		from public.gender g
		join public.status s on g.status_uuid = s.uuid
	)
	`

	// param = append(param, tenantUUID)
	queryBuilder := ""

	//* build query by payload data
	// search
	queryBuilder += `(lower(name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		` or lower(abbr_name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		`)`
	if payload.Search != nil && len(*payload.Search) > 0 {
		param = append(param, "%"+*payload.Search+"%")
	} else {
		param = append(param, "%"+""+"%")
	}

	// filter
	if payload.Filter != nil {
		if (*payload.Filter)["status"] != nil {
			queryBuilder += ` and lower(status) = lower($` + strconv.Itoa(len(param)+1) + `)`
			param = append(param, (*payload.Filter)["status"].(string))
		}
		if (*payload.Filter)["uuid"] != nil {
			queryBuilder += ` and datas.uuid = $` + strconv.Itoa(len(param)+1)
			param = append(param, (*payload.Filter)["uuid"].(string))
		}
	}

	// run count first to get data statistic
	queryCount := query + "SELECT COUNT(*) FROM datas WHERE " + queryBuilder
	count, err := db.GetSingleDataByQuery[model.CountResult](queryCount, param...)
	if err != nil {
		return nil, nil, err
	}

	// order by
	if payload.SortBy != nil {
		queryBuilder += ` ORDER BY `
		for i, sortBy := range *payload.SortBy {
			for key, value := range sortBy {
				if strings.ToLower(value.(string)) == "asc" || strings.ToLower(value.(string)) == "desc" {
					queryBuilder += key + ` ` + value.(string)
					if i+1 < len(*payload.SortBy) {
						queryBuilder += `, `
					}
				}
			}
		}
	}

	// limit
	if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, *payload.RowPerPage)
	} else {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_ROW_PER_PAGES)
	}

	// offset
	if payload.Page != nil && *payload.Page != 0 {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
			param = append(param, *payload.Page**payload.RowPerPage-*payload.RowPerPage)
		} else {
			param = append(param, *payload.Page*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
		}
	} else {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_PAGES*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
	}

	if len(queryBuilder) > 0 {
		query += `SELECT * FROM datas WHERE ` + queryBuilder
	}

	selectedData, err := db.GetMultipleDataByQuery[model.ReadGenderResult](query, param...)
	if err != nil {
		return nil, nil, err
	}

	dataStat := CalculateDataStatisticResult(count, payload, len(*selectedData))

	return *selectedData, &dataStat, nil
}

func SearchTitle(tenantUUID string, payload model.SearchPayload) ([]model.ReadResultTitle, *model.DataStatistics, error) {
	var param []interface{}

	//* base query
	query := `
	with datas as(
		select distinct on (t.abbr_name)
			t.uuid
			,t.name
			,t.abbr_name 
			,t.is_prefix
			,t.sequence
			,s.name as status
		from public.title t
		join public.status s on t.status_uuid = s.uuid
	)
	`

	// param = append(param, tenantUUID)
	queryBuilder := ""

	//* build query by payload data
	// search
	queryBuilder += `(lower(name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		` or lower(abbr_name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		`)`
	if payload.Search != nil && len(*payload.Search) > 0 {
		param = append(param, "%"+*payload.Search+"%")
	} else {
		param = append(param, "%"+""+"%")
	}

	// filter
	if payload.Filter != nil {
		if (*payload.Filter)["status"] != nil {
			queryBuilder += ` and lower(status) = lower($` + strconv.Itoa(len(param)+1) + `)`
			param = append(param, (*payload.Filter)["status"].(string))
		}
		if (*payload.Filter)["uuid"] != nil {
			queryBuilder += ` and datas.uuid = $` + strconv.Itoa(len(param)+1)
			param = append(param, (*payload.Filter)["uuid"].(string))
		}
	}

	// run count first to get data statistic
	queryCount := query + "SELECT COUNT(*) FROM datas WHERE " + queryBuilder
	count, err := db.GetSingleDataByQuery[model.CountResult](queryCount, param...)
	if err != nil {
		return nil, nil, err
	}

	// order by
	if payload.SortBy != nil {
		queryBuilder += ` ORDER BY `
		for i, sortBy := range *payload.SortBy {
			for key, value := range sortBy {
				if strings.ToLower(value.(string)) == "asc" || strings.ToLower(value.(string)) == "desc" {
					queryBuilder += key + ` ` + value.(string)
					if i+1 < len(*payload.SortBy) {
						queryBuilder += `, `
					}
				}
			}
		}
	}

	// limit
	if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, *payload.RowPerPage)
	} else {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_ROW_PER_PAGES)
	}

	// offset
	if payload.Page != nil && *payload.Page != 0 {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
			param = append(param, *payload.Page**payload.RowPerPage-*payload.RowPerPage)
		} else {
			param = append(param, *payload.Page*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
		}
	} else {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_PAGES*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
	}

	if len(queryBuilder) > 0 {
		query += `SELECT * FROM datas WHERE ` + queryBuilder
	}

	selectedData, err := db.GetMultipleDataByQuery[model.ReadResultTitle](query, param...)
	if err != nil {
		return nil, nil, err
	}

	dataStat := CalculateDataStatisticResult(count, payload, len(*selectedData))

	return *selectedData, &dataStat, nil
}

func SearchEducationLevel(tenantUUID string, payload model.SearchPayload) ([]model.ReadResultEducationLevel, *model.DataStatistics, error) {
	var param []interface{}

	//* base query
	query := `
	with datas as(
		select 
			el.uuid
			,el.code
			,el.name
			,el.level_order
			,el.equivalent_level 
			,s."name" as status
		from public.education_level el 
		join public.status s on el.status_uuid = s.uuid
	)
	`

	// param = append(param, tenantUUID)
	queryBuilder := ""

	//* build query by payload data
	// search
	queryBuilder += `(lower(name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		` or lower(code) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		`)`
	if payload.Search != nil && len(*payload.Search) > 0 {
		param = append(param, "%"+*payload.Search+"%")
	} else {
		param = append(param, "%"+""+"%")
	}

	// filter
	if payload.Filter != nil {
		if (*payload.Filter)["status"] != nil {
			queryBuilder += ` and lower(status) = lower($` + strconv.Itoa(len(param)+1) + `)`
			param = append(param, (*payload.Filter)["status"].(string))
		}
		if (*payload.Filter)["uuid"] != nil {
			queryBuilder += ` and datas.uuid = $` + strconv.Itoa(len(param)+1)
			param = append(param, (*payload.Filter)["uuid"].(string))
		}
	}

	// run count first to get data statistic
	queryCount := query + "SELECT COUNT(*) FROM datas WHERE " + queryBuilder
	count, err := db.GetSingleDataByQuery[model.CountResult](queryCount, param...)
	if err != nil {
		return nil, nil, err
	}

	// order by
	if payload.SortBy != nil {
		queryBuilder += ` ORDER BY `
		for i, sortBy := range *payload.SortBy {
			for key, value := range sortBy {
				if strings.ToLower(value.(string)) == "asc" || strings.ToLower(value.(string)) == "desc" {
					queryBuilder += key + ` ` + value.(string)
					if i+1 < len(*payload.SortBy) {
						queryBuilder += `, `
					}
				}
			}
		}
	}

	// limit
	if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, *payload.RowPerPage)
	} else {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_ROW_PER_PAGES)
	}

	// offset
	if payload.Page != nil && *payload.Page != 0 {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
			param = append(param, *payload.Page**payload.RowPerPage-*payload.RowPerPage)
		} else {
			param = append(param, *payload.Page*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
		}
	} else {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_PAGES*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
	}

	if len(queryBuilder) > 0 {
		query += `SELECT * FROM datas WHERE ` + queryBuilder
	}

	selectedData, err := db.GetMultipleDataByQuery[model.ReadResultEducationLevel](query, param...)
	if err != nil {
		return nil, nil, err
	}

	dataStat := CalculateDataStatisticResult(count, payload, len(*selectedData))

	return *selectedData, &dataStat, nil
}

func SearchPosition(tenantUUID string, payload model.SearchPayload) ([]model.ReadResultPosition, *model.DataStatistics, error) {
	var param []interface{}

	//* base query
	query := `
	with datas as(
		select 
			p.uuid
			,p.name
			,p.abbr_name
			,p.is_staff
			,s."name" as status
		from public."position" p
		join public.status s on p.status_uuid = s.uuid
	)
	`

	// param = append(param, tenantUUID)
	queryBuilder := ""

	//* build query by payload data
	// search
	queryBuilder += `(lower(name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		` or lower(abbr_name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		`)`
	if payload.Search != nil && len(*payload.Search) > 0 {
		param = append(param, "%"+*payload.Search+"%")
	} else {
		param = append(param, "%"+""+"%")
	}

	// filter
	if payload.Filter != nil {
		if (*payload.Filter)["status"] != nil {
			queryBuilder += ` and lower(status) = lower($` + strconv.Itoa(len(param)+1) + `)`
			param = append(param, (*payload.Filter)["status"].(string))
		}
		if (*payload.Filter)["is_staff"] != nil {
			queryBuilder += ` and datas.is_staff = $` + strconv.Itoa(len(param)+1)
			param = append(param, (*payload.Filter)["is_staff"].(bool))
		}
	}

	// run count first to get data statistic
	queryCount := query + "SELECT COUNT(*) FROM datas WHERE " + queryBuilder
	count, err := db.GetSingleDataByQuery[model.CountResult](queryCount, param...)
	if err != nil {
		return nil, nil, err
	}

	// order by
	if payload.SortBy != nil {
		queryBuilder += ` ORDER BY `
		for i, sortBy := range *payload.SortBy {
			for key, value := range sortBy {
				if strings.ToLower(value.(string)) == "asc" || strings.ToLower(value.(string)) == "desc" {
					queryBuilder += key + ` ` + value.(string)
					if i+1 < len(*payload.SortBy) {
						queryBuilder += `, `
					}
				}
			}
		}
	}

	// limit
	if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, *payload.RowPerPage)
	} else {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_ROW_PER_PAGES)
	}

	// offset
	if payload.Page != nil && *payload.Page != 0 {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
			param = append(param, *payload.Page**payload.RowPerPage-*payload.RowPerPage)
		} else {
			param = append(param, *payload.Page*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
		}
	} else {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_PAGES*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
	}

	if len(queryBuilder) > 0 {
		query += `SELECT * FROM datas WHERE ` + queryBuilder
	}

	selectedData, err := db.GetMultipleDataByQuery[model.ReadResultPosition](query, param...)
	if err != nil {
		return nil, nil, err
	}

	dataStat := CalculateDataStatisticResult(count, payload, len(*selectedData))

	return *selectedData, &dataStat, nil
}

func SearchEmployeeStatus(tenantUUID string, payload model.SearchPayload) ([]model.ReadResultEmployeeStatus, *model.DataStatistics, error) {
	var param []interface{}

	//* base query
	query := `
	with datas as(
		select 
			s.uuid
			,s."name" 
			,s.abbr_name 
		from public.status s 
	)
	`

	// param = append(param, tenantUUID)
	queryBuilder := ""

	//* build query by payload data
	// search
	queryBuilder += `(lower(name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		` or lower(abbr_name) LIKE lower($` + strconv.Itoa(len(param)+1) + `) ` +
		`)`
	if payload.Search != nil && len(*payload.Search) > 0 {
		param = append(param, "%"+*payload.Search+"%")
	} else {
		param = append(param, "%"+""+"%")
	}

	// filter
	if payload.Filter != nil {
		if (*payload.Filter)["is_staff"] != nil && (*payload.Filter)["is_staff"].(bool) {
			queryBuilder += ` and lower(abbr_name) like lower($` + strconv.Itoa(len(param)+1) + `)`
			param = append(param, `%_staff_%`)
		} else {
			queryBuilder += ` and lower(abbr_name) like lower(%_staff_%)`
			queryBuilder += ` and lower(abbr_name) like lower($` + strconv.Itoa(len(param)+1) + `)`
			param = append(param, `%_employee_%`)
		}
	}

	// run count first to get data statistic
	queryCount := query + "SELECT COUNT(*) FROM datas WHERE " + queryBuilder
	count, err := db.GetSingleDataByQuery[model.CountResult](queryCount, param...)
	if err != nil {
		return nil, nil, err
	}

	// order by
	if payload.SortBy != nil {
		queryBuilder += ` ORDER BY `
		for i, sortBy := range *payload.SortBy {
			for key, value := range sortBy {
				if strings.ToLower(value.(string)) == "asc" || strings.ToLower(value.(string)) == "desc" {
					queryBuilder += key + ` ` + value.(string)
					if i+1 < len(*payload.SortBy) {
						queryBuilder += `, `
					}
				}
			}
		}
	}

	// limit
	if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, *payload.RowPerPage)
	} else {
		queryBuilder += ` LIMIT $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_ROW_PER_PAGES)
	}

	// offset
	if payload.Page != nil && *payload.Page != 0 {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		if payload.RowPerPage != nil && *payload.RowPerPage != 0 {
			param = append(param, *payload.Page**payload.RowPerPage-*payload.RowPerPage)
		} else {
			param = append(param, *payload.Page*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
		}
	} else {
		queryBuilder += ` OFFSET $` + strconv.Itoa(len(param)+1)
		param = append(param, DEFAULT_PAGES*DEFAULT_ROW_PER_PAGES-DEFAULT_ROW_PER_PAGES)
	}

	if len(queryBuilder) > 0 {
		query += `SELECT * FROM datas WHERE ` + queryBuilder
	}

	selectedData, err := db.GetMultipleDataByQuery[model.ReadResultEmployeeStatus](query, param...)
	if err != nil {
		return nil, nil, err
	}

	dataStat := CalculateDataStatisticResult(count, payload, len(*selectedData))

	return *selectedData, &dataStat, nil
}
