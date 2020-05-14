package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type employees struct {
	deptName  string `db:"dept_name"`
	title     string `db:"title"`
	firstName string `db:"first_name"`
	lastName  string `db:"last_name"`
	salary    int    `db:"salary"`
	hireDate  string `db:"hire_date"`
	yearsWork int    `db:"years_work"`
}

type departments struct {
	deptName      string `db:"dept_name"`
	countEmployee int    `db:"count_employee"`
	sumSalary     int    `db:"sum_salary"`
}

const firstQuery string = `
	SELECT t.title, e.first_name, e.last_name, s.salary

	FROM 	employees as e
	LEFT JOIN	titles as t ON t.emp_no = e.emp_no
	LEFT JOIN    salaries as s ON s.emp_no = e.emp_no
	LEFT JOIN    dept_manager as d ON d.emp_no = e.emp_no
	
	WHERE s.to_date > CURRENT_DATE() AND d.to_date > CURRENT_DATE() AND t.title = 'Manager'
	`
const secondQuery string = `
	SELECT  d.dept_name,
        t.title, 
        e.first_name, 
        e.last_name, 				
        e.hire_date,
        (YEAR(CURRENT_DATE)-YEAR(e.hire_date)) - (RIGHT(CURRENT_DATE,5)<RIGHT(e.hire_date,5)) AS years_work
        

	FROM 	employees as e
	LEFT JOIN	titles as t ON t.emp_no = e.emp_no
	LEFT JOIN   dept_emp as de ON de.emp_no = e.emp_no
	LEFT JOIN	departments as d ON d.dept_no = de.dept_no

	WHERE MONTH(e.hire_date) = MONTH(CURRENT_DATE) AND de.to_date > CURRENT_DATE() AND t.to_date > CURRENT_DATE
	LIMIT 10
	`
const thirdQuery string = `
	SELECT d.dept_name, 
	   COUNT(de.emp_no) AS 'count_employee',
       SUM(s.salary) AS 'sum_salary'

	FROM departments as d
	LEFT JOIN dept_emp AS de ON de.dept_no = d.dept_no
	LEFT JOIN employees AS e ON e.emp_no = de.emp_no
	LEFT JOIN salaries AS s ON s.emp_no = e.emp_no

	WHERE de.to_date > CURRENT_DATE AND s.to_date > CURRENT_DATE
	GROUP By d.dept_name
	`

func main() {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/employees")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	findCurrentManagers(firstQuery, db)

	findAllCurrentEmployees(secondQuery, db)

	findDepartaments(thirdQuery, db)

}

func findCurrentManagers(query string, db *sql.DB) {

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	queryStruct := []employees{}

	for rows.Next() {
		p := employees{}
		err := rows.Scan(&p.title, &p.firstName, &p.lastName, &p.salary)
		if err != nil {
			fmt.Println(err)
			continue
		}
		queryStruct = append(queryStruct, p)
	}

	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("| Title | First Name | Last Name | Current Salary |")
	for _, p := range queryStruct {
		fmt.Println(p.title, p.firstName, p.lastName, p.salary)
	}
	fmt.Println("---------------------------------------------------------------------")
}

func findAllCurrentEmployees(query string, db *sql.DB) {

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	queryStruct := []employees{}

	for rows.Next() {
		p := employees{}
		err := rows.Scan(&p.deptName, &p.title, &p.firstName, &p.lastName, &p.hireDate, &p.yearsWork)
		if err != nil {
			fmt.Println(err)
			continue
		}
		queryStruct = append(queryStruct, p)
	}

	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("| Departaments | Title | First Name | Last Name | Hire Date | Years Work |")
	for _, p := range queryStruct {
		fmt.Println(p.deptName, p.title, p.firstName, p.lastName, p.hireDate, p.yearsWork)
	}
	fmt.Println("---------------------------------------------------------------------")
}

func findDepartaments(query string, db *sql.DB) {

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	queryStruct := []departments{}

	for rows.Next() {
		p := departments{}
		err := rows.Scan(&p.deptName, &p.countEmployee, &p.sumSalary)
		if err != nil {
			fmt.Println(err)
			continue
		}
		queryStruct = append(queryStruct, p)
	}

	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("| Departaments | Count employees| Sum salary |")
	for _, p := range queryStruct {
		fmt.Println(p.deptName, p.countEmployee, p.sumSalary)
	}
	fmt.Println("---------------------------------------------------------------------")

}
