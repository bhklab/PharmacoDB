# Cell Lines

```
GET /cell_lines
```

## Description

This method returns a list of cell lines.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- Meta information is included in response headers by default. Use `include` parameter to add info to response body.

## Sources

- http://pharmacodb.pmgenomics.ca/cell_lines

## Parameters

```
GET /cell_lines
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **all** | filter | *boolean* | no | false | Return all items in resource |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |

## Output Formats

- JSON

## Examples

```
GET /cell_lines
```

- https://api.pharmacodb.com/v1/cell_lines

## Output

**JSON**, with metadata included in body.

```
{
    "data": [
        {
            "id": 1,
            "name": "105KC"
        },
        {
            "id": 2,
            "name": "1321N1"
        },
        {
            "id": 3,
            "name": "143B"
        },
        {
            "id": 4,
            "name": "184A1"
        },
        {
            "id": 5,
            "name": "184B5"
        },
        {
            "id": 6,
            "name": "2004"
        },
        {
            "id": 7,
            "name": "201T"
        },
        {
            "id": 8,
            "name": "21MT1"
        },
        {
            "id": 9,
            "name": "21MT2"
        },
        {
            "id": 10,
            "name": "21NT"
        },
        {
            "id": 11,
            "name": "21PT"
        },
        {
            "id": 12,
            "name": "22RV1"
        },
        {
            "id": 13,
            "name": "23132-87"
        },
        {
            "id": 14,
            "name": "253J"
        },
        {
            "id": 15,
            "name": "253J-BV"
        },
        {
            "id": 16,
            "name": "42-MG-BA"
        },
        {
            "id": 17,
            "name": "451Lu"
        },
        {
            "id": 18,
            "name": "501A"
        },
        {
            "id": 19,
            "name": "537 MEL"
        },
        {
            "id": 20,
            "name": "5637"
        },
        {
            "id": 21,
            "name": "59M"
        },
        {
            "id": 22,
            "name": "600MPE"
        },
        {
            "id": 23,
            "name": "624 mel"
        },
        {
            "id": 24,
            "name": "639-V"
        },
        {
            "id": 25,
            "name": "647-V"
        },
        {
            "id": 26,
            "name": "697"
        },
        {
            "id": 27,
            "name": "769-P"
        },
        {
            "id": 28,
            "name": "786-0"
        },
        {
            "id": 29,
            "name": "8-MG-BA"
        },
        {
            "id": 30,
            "name": "8305C"
        }
    ],
    "metadata": {
        "last_page": 57,
        "page": 1,
        "per_page": 30,
        "total": 1691
    }
}
```
