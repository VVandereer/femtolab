# femtolab
**femtolab** — программа для управления установкой фемтосекундного лазера

## Structure 
```plaintext
/ (femtolab)
├── control /              # Интерфейсы управления
│   ├── cli.go             # CLI интерфейс
│   ├── gui.go             # (maybe) графический интерфейс
│   └── api_server.go      # (maybe) Web API
├── main.go                # Ядро приложения
├── logger.go              # Модуль логгера
├── devices_list.go        # HAL интерфейсы устройств и их описание
├── experiment.go          # Модуль для обработки сценария эксперимента
├── devices/               # Модули устройств
│   ├── stepper_motor.go   # Модуль шагового мотора
│   ├── laser_generator.go # Модуль лазерного генератора
│   ├── oscilloscope.go    # Модуль осциллографа
│   ├── lockin_sr830.go    # Модуль LOCKIN SR830
│   └── camera.go          # Модуль камеры
├── scripts/               # Примеры и сценарии экспериментов
│   └── example_experiment.yaml
├── config.yaml            # Конфигурационные файлы устройства и системы
```
