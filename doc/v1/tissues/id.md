# Tissue

```
GET /tissues/{id}
```

## Description

This method returns a single tissue.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- A `404` error is returned if an item is not found.

## Sources

- http://pharmacodb.ca/tissues

## Parameters

```
GET /tissues/{id}
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | tissue_id | Define whether `id = tissue_id` or `id = tissue_name` |

## Output Formats

- JSON

## Examples

```
GET /tissues/{id}
```

- http://api.pharmacodb.ca/v1/tissues/1
- http://api.pharmacodb.ca/v1/tissues/breast?type=name
- http://api.pharmacodb.ca/v1/tissues/7?indent=true

## Output

**JSON**, using the tissue **breast**.

```
{
    "id": 7,
    "name": "breast",
    "annotations": [
        {
            "name": "breast",
            "datasets": [
                "CCLE",
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ]
        },
        {
            "name": "Breast",
            "datasets": [
                "GDSC1000",
                "gCSI"
            ]
        },
        {
            "name": "breast ",
            "datasets": [
                "GDSC1000"
            ]
        }
    ]
}
```
