# Endpoint

```
GET /{endpoint}
```

## Description

This is a general endpoint doc for endpoints that have the same documentation with a single parameter (`indent`).

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Parameters

```
GET /{endpoint}
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **indent** | input | *boolean* | no | false | Add indentation to response |

## Output Formats

- JSON

## Examples

```
GET /{endpoint}
```

- https://api.pharmacodb.com/v1/intersections/1/895/526
- https://api.pharmacodb.com/v1/intersections/2/895/1?indent=true
- https://api.pharmacodb.com/v1/stats/tissues/cell_lines
- https://api.pharmacodb.com/v1/stats/datasets/cell_lines
- https://api.pharmacodb.com/v1/stats/datasets/tissues
- https://api.pharmacodb.com/v1/stats/datasets/compounds
- https://api.pharmacodb.com/v1/stats/datasets/experiments
