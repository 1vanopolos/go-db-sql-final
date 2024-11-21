package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	res, err := s.db.Exec("INSERT INTO parcel (address, client, created_at, status) VALUES (:address, :client, :created_at, :status)",
		sql.Named("address", p.Address),
		sql.Named("client", p.Client),
		sql.Named("created_at", p.CreatedAt),
		sql.Named("status", p.Status))
	if err != nil {
		return 0, err
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

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	row := s.db.QueryRow("SELECT number, address, client, created_at, status FROM parcel WHERE number = :n", sql.Named("n", number))
	err := row.Scan(&p.Number, &p.Address, &p.Client, &p.CreatedAt, &p.Status)

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	rows, err := s.db.Query("SELECT number, address, client, created_at, status FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	for rows.Next() {

		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Address, &p.Client, &p.CreatedAt, &p.Status)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
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

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	p, err := s.Get(number)
	if err != nil {
		return err
	}
	// менять адрес можно только если значение статуса registered
	if p.Status != ParcelStatusRegistered {
		return err
	}
	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) Delete(number int) error {
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
	return err

}
