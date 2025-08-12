Изначально проект должен был быть на основе таких таблиц:
```
Date,DayNumber,KeyTask,Category,StressBefore_0_10,Started_Y_N,StartTime_HH_MM,ActiveDuration_min,ContinuedAfter10min_Y_N,StressAfter_0_10,StressReduction,Distractions_min,BlocksCompleted,PomodoroCount,LightExposure_min,Energy_0_10,Mood_0_10,Notes
2025-08-12,1,Написать отчет по проекту,работа,8,Yes,09:10,30,Yes,4,=E2-J2,5,2,1,15,5,4,Отвлекся на почту 5 мин
,,,,,,,,,,,,,,,,,Заполните строку: запуск таймера 10 минут и внесение данных
```

```
Date,Bedtime_HH_MM,WakeTime_HH_MM,SleepLatency_min,NightAwakenings_count,TotalSleepHours,SleepQuality_0_10,DaytimeSleepiness_0_10,CaffeineAfterNoon_Y_N,ScreenUseBeforeBed_min,EveningFreeTime_min,Notes
2025-08-11,00:30,08:00,30,2,7.5,6,5,Yes,90,30,Запланировать свет утром 15 минут
,,,,,,,,,,,,,,
```


Directories
```
.
├── cmd
│   └── api
├── .gitignore
├── go.mod
├── internal
│   ├── application
│   │   ├── commands
│   │   ├── handlers
│   │   ├── queries
│   │   └── services
│   ├── domain
│   │   ├── aggregates
│   │   ├── entities
│   │   │   ├── sleep_entry.go
│   │   │   ├── task_entry.go
│   │   │   └── task_entry_test.go
│   │   ├── events
│   │   │   ├── domain_event.go
│   │   │   └── task_events.go
│   │   ├── repositories
│   │   │   └── task_repository.go
│   │   └── valueobjects
│   │       ├── levels.go
│   │       └── levels_test.go
│   ├── infrastructure
│   │   ├── config
│   │   ├── http
│   │   └── persistence
│   └── interfaces
│       ├── dto
│       └── rest
├── migrations
├── pkg
│   ├── errors
│   │   └── domain_error.go
│   ├── logger
│   └── utils
├── README.md
└── web
```