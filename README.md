# femtolab
**femtolab** — программа для управления установкой фемтосекундного лазера

## Structure 
```plaintext
/ (femtolab)
├── control /                 # Control interfaces
│ ├── cli.go                  # CLI interface
│ ├── gui.go                  # (maybe) graphical interface
│ └── api_server.go           # (maybe) Web API
├── main.go                   # Application core
├── logger.go                 # Logger module
├── devices_list.go           # HAL device interfaces and their description
├── experiment.go             # Module for processing the experiment scenario
├── devices/                  # Device modules
│ ├── stepper_motor.go        # Stepper motor module
│ ├── laser_generator.go      # Laser Generator Module
│ ├── oscilloscope.go         # Oscilloscope Module
│ ├── lockin_sr830.go         # LOCKIN SR830 Module
│ └── camera.go               # Camera Module
├── scripts/                  # Examples and Experiment Scenarios
│ └── example_experiment.yaml
├── config.yaml               # Device and System Configuration Files
```
