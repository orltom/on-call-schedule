[![CI](https://github.com/orltom/on-call-schedule/actions/workflows/github-actions-ci.yaml/badge.svg)](https://github.com/orltom/on-call-schedule/actions/workflows/github-actions-ci.yaml)
[![License](https://img.shields.io/github/license/orltom/on-call-schedule)](/LICENSE)

# On-Call Schedule
This tool is a flexible command-line tool designed to simplify the creation and management of on-call schedules 
for on-call teams. It addresses common limitations in existing incident management tools, such as lack of calendar 
syncing, inadequate handling of holidays, and inflexible scheduling rules.

# ✨ Current Features
- Custom Rule-Based On-Call Scheduling: Create on-call schedules based on your own rules or use default rules, such as holiday checkers.
- Multiple Export Formats: Generate on-call schedules in various formats, including CSV, iCalender and JSON

# 🛠️ Planned Features
- Automatically generate on-call schedules by syncing with a shared group calendar.
- Export schedules directly to popular incident management tools.


# Installation
To install the latest version of ocsctl:
```shell
go install github.com/orltom/on-call-schedule@latest
```

Or clone the repository and build manually:

```shell
git clone https://github.com/orltom/on-call-schedule.git
go build -o ocsctl ./cmd
```

## Usage
### Create on-call schedule plan
Here is an example of how to create an on-call duty plan
```shell
cat > team.json << EOL
{
  "employees": [
    {"id": "joe@example.com", "name": "Joe"},
    {"id": "jan@example.com", "name": "Jan", "vacationDays": ["2024-01-06","2024-01-07"]},
    {"id": "lee@example.com", "name": "Lee"},
    {"id": "eva@example.com", "name": "Eva"}
  ]
}
EOL

ocsctl create \
        --start "2024-01-01 00:00:00" \
        --end "2024-03-29 00:00:00" \
        --duration 168 \
        --team-file team.json \
        --primary-rules=vacation,minimumfourshiftgap
        --secondary-rules=vacation,minimumtwoshiftgap
        --output table
```
## Help Command
```shell
ocsctl [command] -h
```


## Contributing
Contributions are welcome in any form, be it code, logic, documentation, examples, requests, bug reports, ideas or
anything else that will help this project move forward.

## License
This project is licensed under the MIT License. See the LICENSE file for more details.
