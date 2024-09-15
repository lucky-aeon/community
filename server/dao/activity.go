package dao

import (
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

type ActivityDAO struct {
}

// 课程日活查询
func (d ActivityDAO) GetCourseDailyActivity() ([]model.CourseActivity, error) {
	var results []model.CourseActivity
	query := `
        SELECT c.title AS course_title, IFNULL(COUNT(DISTINCT ol.user_id), 0) AS daily_active_users
        FROM courses c
        JOIN courses_sections cs ON cs.course_id = c.id
        LEFT JOIN community.oper_logs ol ON cs.id = SUBSTRING_INDEX(ol.request_info, '/', -1)
        WHERE ol.request_info LIKE '/community/courses/section/%'
        AND DATE(ol.created_at) = CURDATE()
        GROUP BY c.title;
    `

	if err := mysql.GetInstance().Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		return make([]model.CourseActivity, 0), nil
	}
	return results, nil
}

// 课程周活查询
func (d *ActivityDAO) GetCourseWeeklyActivity() ([]model.CourseWeekActivity, error) {
	var results []model.CourseWeekActivity
	query := `
      SELECT
    c.title AS course_title,
    IFNULL(COUNT(DISTINCT ol.user_id), 0) AS weekly_active_users
FROM
    courses c
        JOIN
    courses_sections cs ON cs.course_id = c.id
        LEFT JOIN
    community.oper_logs ol ON cs.id = SUBSTRING_INDEX(ol.request_info, '/', -1)
WHERE
    ol.request_info LIKE '/community/courses/section/%'
  AND ol.created_at BETWEEN DATE_SUB(CURDATE(), INTERVAL 7 DAY) AND CURDATE()  -- 从今天往前推7天
GROUP BY
    c.title;

    `
	if err := mysql.GetInstance().Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		return make([]model.CourseWeekActivity, 0), nil
	}
	return results, nil
}

// 课程月活查询
func (d *ActivityDAO) GetCourseMonthlyActivity() ([]model.CourseMonthActivity, error) {
	var results []model.CourseMonthActivity
	query := `
        SELECT c.title AS course_title, IFNULL(COUNT(DISTINCT ol.user_id), 0) AS monthly_active_users
        FROM courses c
        JOIN courses_sections cs ON cs.course_id = c.id
        LEFT JOIN community.oper_logs ol ON cs.id = SUBSTRING_INDEX(ol.request_info, '/', -1)
        WHERE ol.request_info LIKE '/community/courses/section/%'
        AND YEAR(ol.created_at) = YEAR(CURDATE())
        AND MONTH(ol.created_at) = MONTH(CURDATE())
        GROUP BY c.title;
    `
	if err := mysql.GetInstance().Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		return make([]model.CourseMonthActivity, 0), nil
	}
	return results, nil
}

// 一周内每日活跃度趋势
func (d *ActivityDAO) GetCourseWeeklyTrend() ([]model.CourseActivity, error) {
	var results []model.CourseActivity
	query := `
SELECT 
    c.title AS course_title,
    dates.date AS activity_date,
    IFNULL(COUNT(DISTINCT ol.user_id), 0) AS daily_active_users
FROM 
    (
        SELECT CURDATE() AS date
        UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 1 DAY)
        UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 2 DAY)
        UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 3 DAY)
        UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 4 DAY)
        UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 5 DAY)
        UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 6 DAY)
    ) AS dates
JOIN 
    courses c  -- 用 LEFT JOIN 避免生成过多的笛卡尔积
LEFT JOIN 
    courses_sections cs ON cs.course_id = c.id
LEFT JOIN 
    community.oper_logs ol ON cs.id = SUBSTRING_INDEX(ol.request_info, '/', -1)
    AND DATE(ol.created_at) = dates.date
WHERE 
    ol.request_info LIKE '/community/courses/section/%'
GROUP BY 
    c.title, dates.date
ORDER BY 
    c.title, dates.date;

    `
	if err := mysql.GetInstance().Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		return make([]model.CourseActivity, 0), nil
	}
	return results, nil
}

// 查询日活
func (d *ActivityDAO) GetUserDailyActiveUsers() (model.ActivityData, error) {
	var result model.ActivityData
	query := `
        SELECT COUNT(DISTINCT user_id) AS active_users
        FROM community.oper_logs
        WHERE DATE(created_at) = CURDATE();
    `
	if err := mysql.GetInstance().Raw(query).Scan(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

// 查询周活
func (d *ActivityDAO) GetUserWeeklyActiveUsers() (model.ActivityData, error) {
	var result model.ActivityData
	query := `
       SELECT 
    COUNT(DISTINCT user_id) AS active_users
FROM 
    community.oper_logs
WHERE 
    created_at BETWEEN DATE_SUB(CURDATE(), INTERVAL 7 DAY) AND CURDATE();

    `
	if err := mysql.GetInstance().Raw(query).Scan(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

// 查询月活
func (d *ActivityDAO) GetUserMonthlyActiveUsers() (model.ActivityData, error) {
	var result model.ActivityData
	query := `
        SELECT COUNT(DISTINCT user_id) AS active_users
        FROM community.oper_logs
        WHERE YEAR(created_at) = YEAR(CURDATE())
        AND MONTH(created_at) = MONTH(CURDATE());
    `
	if err := mysql.GetInstance().Raw(query).Scan(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

// 查询人数折线图
func (d *ActivityDAO) GetUserActivityLineData() ([]model.ActivityLine, error) {
	var results []model.ActivityLine
	query := `
        SELECT dates.date, IFNULL(COUNT(DISTINCT ol.user_id), 0) AS active_users
        FROM (
            SELECT CURDATE() AS date
            UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 1 DAY)
            UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 2 DAY)
            UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 3 DAY)
            UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 4 DAY)
            UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 5 DAY)
            UNION ALL SELECT DATE_SUB(CURDATE(), INTERVAL 6 DAY)
        ) AS dates
        LEFT JOIN community.oper_logs ol
        ON DATE(ol.created_at) = dates.date
        GROUP BY dates.date
        ORDER BY dates.date;
    `
	if err := mysql.GetInstance().Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
