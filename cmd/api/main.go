package main

import (
	"car-service/internal/application/commands"
	"car-service/internal/application/services"
	"car-service/internal/domain/repositories"
	"car-service/internal/infrastructure/api/handlers"
	api "car-service/internal/infrastructure/api/mediator"
	"car-service/internal/infrastructure/api/server"
	gormrepo "car-service/internal/infrastructure/persistence/gorm"
	"car-service/internal/infrastructure/persistence/migrations"
	"car-service/pkg/config"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Cargar configuración de ambiente
	env, err := config.NewEnvironment()
	if err != nil {
		return err
	}

	// Configuración de la base de datos
	db, err := setupDatabase(env)
	if err != nil {
		return err
	}

	// Inicializar repositorios
	var carRepo repositories.CarRepository = gormrepo.NewCarRepository(db)
	carService := services.NewCarService(carRepo)

	mediator := api.NewMediator()

	newCarCommand := commands.NewNewCarCommand(carService)
	mediator.Register("CreateCar", newCarCommand)

	carHandler := handlers.NewCarHandler(mediator)

	// Configurar el servidor
	serverCfg := &server.ServerConfig{
		CarHandler: carHandler,
		Port:       env.ServerPort,
	}

	// Crear y configurar el servidor
	srv := server.NewServer(serverCfg)

	// Iniciar el servidor (bloqueante)
	return srv.Start()
}
func setupDatabase(env *config.Environment) (*gorm.DB, error) {
	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(env.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ejecutar migraciones
	if err := migrations.Migrate(db); err != nil {
		return nil, err
	}
	log.Println("Migraciones ejecutadas correctamente")

	return db, nil
}
