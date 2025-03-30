# femtolab
**femtolab** — программа для управления установкой фемтосекундного лазера

## Structure 
```plaintext
/ (femtolab)
├── cmd/                   # CLI интерфейс
│   └── femtot-cli.go      # Точка входа для командной строки
├── core/                  # Ядро приложени
│   ├── main.go            # Само ядро
│   ├── experiment.go      # Проведение сценариев экспериментов
│   └── logger.go          # Логирование действий
├── devices/               # Модули устройств
│   ├── stepper_motor.go   # Модуль шагового мотора
│   ├── laser_generator.go # Модуль лазерного генератора
│   ├── oscilloscope.go    # Модуль осциллографа
│   ├── lockin_sr830.go    # Модуль LOCKIN SR830
│   └── camera.go          # Модуль камеры
├── interface/             # Интерфейсы управления
│   ├── cli.go             # CLI интерфейс
│   ├── gui.go             # (maybe) графический интерфейс
│   └── api_server.go      # (maybe) Web API
├── scripts/               # Примеры и сценарии экспериментов
│   └── example_experiment.yaml
├── config/                # Конфигурационные файлы устройств и системы
│   └── devices.yaml
```
