package main

import (
	"car-service/cmd/api/controllers"
	"car-service/cmd/api/mediator"
	"car-service/cmd/api/server"
	"car-service/internal/application/commands/new_car"
	"car-service/internal/application/commands/update_car" // Added for UpdateCar
	"car-service/internal/application/queries/get_car_by_id" // Added for GetCarByID
	"car-service/internal/application/queries/get_cars"
	"car-service/internal/application/services"
	"car-service/internal/domain/repositories"
	gormrepo "car-service/internal/infrastructure/gorm"
	"car-service/internal/infrastructure/migrations"
	"car-service/pkg/config"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Car Service API
// @version 1.0
// @description API para el servicio de autos
// @host localhost:8080
// @BasePath /api
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
	var modelRepo repositories.ModelRepository = gormrepo.NewModelRepository(db)
	var ownerRepo repositories.OwnerRepository = gormrepo.NewOwnerRepository(db)

	// Inicializar servicios
	carService := services.NewCarService(carRepo, modelRepo, ownerRepo)

	// Inicializar mediator y registrar handlers
	m := mediator.NewMediator(db)
	m.RegisterCommand(new_car.Name, new_car.NewNewCarCommandHandler(carService))
	m.RegisterCommand(update_car.Name, update_car.NewUpdateCarCommandHandler(carService)) // Register UpdateCar
	m.RegisterQuery(get_cars.Name, get_cars.NewGetCarsQueryHandler(carService))
	m.RegisterQuery(get_car_by_id.Name, get_car_by_id.NewGetCarByIDQueryHandler(carService)) // Register GetCarByID

	// Inicializar controllers
	carController := controllers.NewCarController(m)

	// Configurar el servidor
	serverCfg := &server.ServerConfig{
		CarController: carController,
		Port:          env.ServerPort,
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
