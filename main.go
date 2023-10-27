package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jstoledano/CECyRDIngestionTool/database"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PGX"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var conteo string
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tramites").Scan(&conteo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Actualmente, existen ", conteo, "trámites en la base de datos")

	directoryPath := "./data"
	dir, err := os.ReadDir(directoryPath)
	if err != nil {
		fmt.Println("Error al leer el directorio:", err)
		return
	}
	fmt.Println("Archivos en el directorio:")
	for _, entry := range dir {
		if entry.IsDir() {
			fmt.Println("Directorio:", entry.Name())
		} else {
			fmt.Println("Archivo:", entry.Name())
			fullPath := filepath.Join(directoryPath, entry.Name())
			f, err := os.Open(fullPath)
			if err != nil {
				log.Fatal(err)
			}
			procesaArchivo(f, conn)
			defer f.Close()
		}
	}
}

func procesaArchivo(f *os.File, conn *pgx.Conn) {
	csvReader := csv.NewReader(f)
	csvReader.Comma = '|'
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	loc, _ := time.LoadLocation("America/Mexico_City")
	layout := "02/01/2006 03:04:05 PM"

	log.Println("Processing", len(data), "records")
	bar := pb.StartNew(len(data))
	start := time.Now()

	for _, row := range data {
		t := database.Tramite{}
		t.Folio = row[0]
		t.Estatus = row[1]
		t.CausaRechazo = row[2]
		t.MovimientoSolicitado = row[3]
		t.MovimientoDefinitivo = row[4]
		t.FechaTramite, _ = time.ParseInLocation(layout, row[5], loc)
		t.FechaRecibidoCecyrd, _ = time.ParseInLocation(layout, row[6], loc)
		t.FechaRegistradoCecyrd, _ = time.ParseInLocation(layout, row[7], loc)
		t.FechaRechazado, _ = time.ParseInLocation(layout, row[8], loc)
		t.FechaCanceladoMovimientoPosterior, _ = time.ParseInLocation(layout, row[9], loc)
		t.FechaAltaPe, _ = time.ParseInLocation(layout, row[10], loc)
		t.FechaAfectacionPadron, _ = time.ParseInLocation(layout, row[11], loc)
		t.FechaActualizacionPe, _ = time.ParseInLocation(layout, row[12], loc)
		t.FechaReincorporacionPe, _ = time.ParseInLocation(layout, row[13], loc)
		t.FechaExitoso, _ = time.ParseInLocation(layout, row[14], loc)
		t.FechaLoteProduccion, _ = time.ParseInLocation(layout, row[15], loc)
		t.FechaListoReimpresion, _ = time.ParseInLocation(layout, row[16], loc)
		t.FechaCpvCreada, _ = time.ParseInLocation(layout, row[17], loc)
		t.FechaCpvRegistradaMac, _ = time.ParseInLocation(layout, row[18], loc)
		t.FechaCpvDisponible, _ = time.ParseInLocation(layout, row[19], loc)
		t.FechaCpvEntregada, _ = time.ParseInLocation(layout, row[20], loc)
		t.FechaAfectacionLn, _ = time.ParseInLocation(layout, row[21], loc)
		t.Distrito, _ = strconv.Atoi(t.Folio[5:6])
		t.Mac = t.Folio[2:8]
		t.TramoDisponible = t.FechaCpvDisponible.Sub(t.FechaTramite)
		t.TramoEntrega = t.FechaCpvEntregada.Sub(t.FechaCpvDisponible)
		t.TramoExitoso = t.FechaExitoso.Sub(t.FechaTramite)

		_, err = conn.Exec(context.Background(),
			`INSERT INTO tramites 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27) ON CONFLICT (folio) DO UPDATE SET 
				estatus = $2, 
				causa_rechazo = $3, 
				movimiento_solicitado = $4, 
				movimiento_definitivo = $5, 
				fecha_tramite = $6, 
				fecha_recibido_cecyrd = $7, 
				fecha_registrado_cecyrd = $8, 
				fecha_rechazado = $9, 
				fecha_cancelado_movimiento_posterior = $10, 
				fecha_alta_pe = $11, 
				fecha_afectacion_padron = $12, 
				fecha_actualizacion_pe = $13, 
				fecha_reincorporacion_pe = $14, 
				fecha_exitoso = $15, 
				fecha_lote_produccion = $16, 
				fecha_listo_reimpresion = $17, 
				fecha_cpv_creada = $18, 
				fecha_cpv_registrada_mac = $19, 
				fecha_cpv_disponible = $20, 
				fecha_cpv_entregada = $21, 
				fecha_afectacion_ln = $22,
				 distrito = $23, 
				 mac = $24, 
				 tramo_disponible = $25, 
				 tramo_entrega = $26, 
				 tramo_exitoso = $27`,
			t.Folio, t.Estatus, t.CausaRechazo, t.MovimientoSolicitado, t.MovimientoDefinitivo,
			t.FechaTramite, t.FechaRecibidoCecyrd, t.FechaRegistradoCecyrd, t.FechaRechazado, t.FechaCanceladoMovimientoPosterior,
			t.FechaAltaPe, t.FechaAfectacionPadron, t.FechaActualizacionPe, t.FechaReincorporacionPe, t.FechaExitoso,
			t.FechaLoteProduccion, t.FechaListoReimpresion, t.FechaCpvCreada, t.FechaCpvRegistradaMac, t.FechaCpvDisponible,
			t.FechaCpvEntregada, t.FechaAfectacionLn, t.Distrito, t.Mac, t.TramoDisponible,
			t.TramoEntrega, t.TramoExitoso)
		if err != nil {
			log.Fatal(err)
		}
		bar.Increment()
	}

	bar.Finish()
	end := time.Now()
	log.Println("Time elapsed:", end.Sub(start))
}
