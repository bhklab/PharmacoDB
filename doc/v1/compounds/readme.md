# Compounds

```
GET /compounds
```

## Description

This method returns a list of compounds.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- Meta information is included in response headers by default. Use `include` parameter to add info to response body.

## Sources

- http://pharmacodb.ca/compounds

## Parameters

```
GET /compounds
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
GET /compounds
```

- http://api.pharmacodb.ca/v1/compounds

## Output

**JSON**, with metadata included in body.

```
{
    "data": [
        {
            "id": 1,
            "name": "(5Z)-7-Oxozeaenol"
        },
        {
            "id": 2,
            "name": "16-beta-bromoandrosterone"
        },
        {
            "id": 3,
            "name": "17-AAG"
        },
        {
            "id": 4,
            "name": "1S,3R-RSL-3"
        },
        {
            "id": 5,
            "name": "3-Cl-AHPC"
        },
        {
            "id": 6,
            "name": "5-FdUR"
        },
        {
            "id": 7,
            "name": "5-FU"
        },
        {
            "id": 8,
            "name": "681640"
        },
        {
            "id": 9,
            "name": "968"
        },
        {
            "id": 10,
            "name": "A-443654"
        },
        {
            "id": 11,
            "name": "A-770041"
        },
        {
            "id": 12,
            "name": "A-804598"
        },
        {
            "id": 13,
            "name": "AA-COCF3"
        },
        {
            "id": 14,
            "name": "abiraterone"
        },
        {
            "id": 15,
            "name": "ABT-199"
        },
        {
            "id": 16,
            "name": "ABT-263"
        },
        {
            "id": 17,
            "name": "ABT-737"
        },
        {
            "id": 18,
            "name": "ABT-888"
        },
        {
            "id": 19,
            "name": "AC220"
        },
        {
            "id": 20,
            "name": "AC55649"
        },
        {
            "id": 21,
            "name": "AEW541"
        },
        {
            "id": 22,
            "name": "AG-014699"
        },
        {
            "id": 23,
            "name": "AG1478"
        },
        {
            "id": 24,
            "name": "AGK-2"
        },
        {
            "id": 25,
            "name": "AICAR"
        },
        {
            "id": 26,
            "name": "AKT inhibitor VIII"
        },
        {
            "id": 27,
            "name": "alisertib"
        },
        {
            "id": 28,
            "name": "alisertib:navitoclax (2:1 mol/mol)"
        },
        {
            "id": 29,
            "name": "alvocidib"
        },
        {
            "id": 30,
            "name": "AM-580"
        }
    ],
    "metadata": {
        "last_page": 26,
        "page": 1,
        "per_page": 30,
        "total": 759
    }
}
```
