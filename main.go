package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Producto struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
	Stock  int     `json:"stock"`
}

var db *sql.DB

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	godotenv.Load()

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error conectando a la BD:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a Supabase:", err)
	}

	log.Println("Conectado a Supabase correctamente")

	crearTabla()

	http.HandleFunc("/productos", corsMiddleware(manejarProductos))
	http.HandleFunc("/productos/", corsMiddleware(manejarProductoPorID))

	log.Println("API corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func crearTabla() {
	query := `CREATE TABLE IF NOT EXISTS productos (
		id SERIAL PRIMARY KEY,
		nombre VARCHAR(100) NOT NULL,
		precio DECIMAL(10,2) NOT NULL,
		stock INT NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creando tabla:", err)
	}
}

func manejarProductos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		obtenerProductos(w, r)
	case "POST":
		crearProducto(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func manejarProductoPorID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/productos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		obtenerProducto(w, r, id)
	case "PUT":
		actualizarProducto(w, r, id)
	case "DELETE":
		eliminarProducto(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func obtenerProductos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nombre, precio, stock FROM productos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var productos []Producto
	for rows.Next() {
		var p Producto
		rows.Scan(&p.ID, &p.Nombre, &p.Precio, &p.Stock)
		productos = append(productos, p)
	}
	json.NewEncoder(w).Encode(productos)
}

func obtenerProducto(w http.ResponseWriter, r *http.Request, id int) {
	var p Producto
	err := db.QueryRow("SELECT id, nombre, precio, stock FROM productos WHERE id=$1", id).
		Scan(&p.ID, &p.Nombre, &p.Precio, &p.Stock)
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func crearProducto(w http.ResponseWriter, r *http.Request) {
	var p Producto
	json.NewDecoder(r.Body).Decode(&p)
	err := db.QueryRow(
		"INSERT INTO productos (nombre, precio, stock) VALUES ($1, $2, $3) RETURNING id",
		p.Nombre, p.Precio, p.Stock,
	).Scan(&p.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func actualizarProducto(w http.ResponseWriter, r *http.Request, id int) {
	var p Producto
	json.NewDecoder(r.Body).Decode(&p)
	p.ID = id
	_, err := db.Exec(
		"UPDATE productos SET nombre=$1, precio=$2, stock=$3 WHERE id=$4",
		p.Nombre, p.Precio, p.Stock, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func eliminarProducto(w http.ResponseWriter, r *http.Request, id int) {
	_, err := db.Exec("DELETE FROM productos WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
