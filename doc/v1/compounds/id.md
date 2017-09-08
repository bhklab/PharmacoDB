# Compound

```
GET /compounds/{id}
```

## Description

This method returns a single compound.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- A `404` error is returned if an item is not found.

## Sources

- http://pharmacodb.pmgenomics.ca/compounds

## Parameters

```
GET /compounds/{id}
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | tissue_id | Define whether `id = compound_id` or `id = compound_name` |

## Output Formats

- JSON

## Examples

```
GET /compounds/{id}
```

- https://api.pharmacodb.com/v1/compounds/1
- https://api.pharmacodb.com/v1/compounds/paclitaxel?type=name
- https://api.pharmacodb.com/v1/compounds/526?indent=true

## Output

**JSON**, using the compound **paclitaxel**.

```
{
    "id": 526,
    "name": "paclitaxel",
    "annotations": [
        {
            "name": "Paclitaxel",
            "datasets": [
                "CCLE",
                "GDSC1000",
                "GRAY",
                "FIMM",
                "UHNBreast"
            ]
        },
        {
            "name": "paclitaxel",
            "datasets": [
                "gCSI",
                "CTRPv2"
            ]
        }
    ]
}
```
