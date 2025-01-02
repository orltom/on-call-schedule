[![CI](https://github.com/orltom/on-call-schedule/actions/workflows/github-actions-ci.yaml/badge.svg)](https://github.com/orltom/on-call-schedule/actions/workflows/github-actions-ci.yaml)
[![License](https://img.shields.io/github/license/orltom/on-call-schedule)](/LICENSE)

# On-Call Schedule
It helps to create or synchronize on-call schedules for team members, allowing users to define custom rules or use 
predefined ones. Important factors, such as absences due to holidays or other time off, are automatically considered.

## Usage
Here is an example of how to create an on-call duty plan
```shell
> cat > demo.json << EOL
{
  "employees": [
    {"id": "joe@example.com", "name": "Joe"},
    {"id": "jan@example.com", "name": "Jan", "vacationDays": ["2024-01-06","2024-01-07"]},
    {"id": "lee@example.com", "name": "Lee"}
  ]
}
EOL

./ocsctl create \
        --start "2024-01-01 00:00:00" \
        --end "2024-03-29 00:00:00" \
        --duration 168 \
        --team-file demo.json \
        --output table
```

## Contributing
Contributions are welcome in any form, be it code, logic, documentation, examples, requests, bug reports, ideas or
anything else that will help this project move forward.
