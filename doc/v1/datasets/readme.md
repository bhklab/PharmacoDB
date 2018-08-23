# Datasets

```
GET /datasets
```

## Description

This method returns a list of datasets.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- Meta information is included in response headers by default. Use `include` parameter to add info to response body.

## Sources

- http://pharmacodb.ca/datasets

## Parameters

```
GET /datasets
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
GET /datasets
```

- http://api.pharmacodb.ca/v1/datasets

## Output

**JSON**, with metadata included in body.

```
{
    "data": [
        {
            "id": 1,
            "name": "CCLE"
        },
        {
            "id": 5,
            "name": "GDSC1000"
        },
        {
            "id": 4,
            "name": "gCSI"
        },
        {
            "id": 6,
            "name": "GRAY"
        },
        {
            "id": 3,
            "name": "FIMM"
        },
        {
            "id": 2,
            "name": "CTRPv2"
        },
        {
            "id": 7,
            "name": "UHNBreast"
        }
    ],
    "metadata": {
        "last_page": 1,
        "page": 1,
        "per_page": 30,
        "total": 7
    }
}
```
