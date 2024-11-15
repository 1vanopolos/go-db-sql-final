package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (number, client, status, address, createdat) VALUES (:number, :client, :status, :address, :createdat)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("client", p.Status),
		sql.Named("client", p.Address),
		sql.Named("client", p.CreatedAt))

	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&number)
	if err != nil {
		fmt.Println(err)
		return Parcel{}, nil
	}
	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		fmt.Println(err)
		return []Parcel{}, nil
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {

		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Address, &p.Client, &p.CreatedAt, &p.Status)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	// Проверка ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel

	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	p, err := s.Get(number)
	if err != nil {
		return err
	}

	// удалять строку можно только если значение статуса registered
	if p.Status != ParcelStatusRegistered {
		return err
	}

	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("status", address),
		sql.Named("number", number))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err

}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	// реализуйте удаление строки из таблицы parcel

	p, err := s.Get(number)
	if err != nil {
		return err
	}

	// удалять строку можно только если значение статуса registered
	if p.Status != ParcelStatusRegistered {
		return err
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))

	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
