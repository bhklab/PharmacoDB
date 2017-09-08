# Cell Line

```
GET /cell_lines/{id}
```

## Description

This method returns a single cell line.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- A `404` error is returned if an item is not found.

## Sources

- http://pharmacodb.pmgenomics.ca/cell_lines

## Parameters

```
GET /cell_lines/{id}
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | cell_id | Define whether `id = cell_id` or `id = cell_name` |

## Output Formats

- JSON

## Examples

```
GET /cell_lines/{id}
```

- https://api.pharmacodb.com/v1/cell_lines/1
- https://api.pharmacodb.com/v1/cell_lines/mcf7?type=name
- https://api.pharmacodb.com/v1/cell_lines/895?indent=true

## Output

**JSON**, using the cell line **MCF7**.

```
{
    "id": 895,
    "name": "MCF7",
    "tissue": {
        "id": 7,
        "name": "breast"
    },
    "annotations": [
        {
            "name": "MCF7",
            "datasets": [
                "CCLE",
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ]
        },
        {
            "name": "MCF-7",
            "datasets": [
                "gCSI"
            ]
        },
        {
            "name": "mcf7",
            "datasets": [
                "UHNBreast"
            ]
        }
    ]
}
```
