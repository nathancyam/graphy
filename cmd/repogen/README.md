# Cypher Repository Generation Tool

This tool should be run with the following command:

```bash
> go run graphy/cmd/repogen/main.go path/to/repo.yaml
```

This will print out a basic repository that attempts to implement the given interface provided in the YAML file. Where the generated code goes depends on your discretion. The YAML file is expected to be formatted with the following structure:

```yaml
name: RoundRepository
implements: graphy/pkg/competition/rounds.Repository
methods:
  - name: FindRoundsByID
    cypher: |-
      MATCH (res:Round) WHERE res.id IN $roundIDs RETURN res LIMIT 10
    output: res
  - name: FindGradeRounds
    cypher: |-
      MATCH (g:Grade)-[:HAS_ROUND]->(r:Round)
      WHERE g.id IN $gradeIDs
      RETURN { id: g.id, items: COLLECT(r) } as out
    output: out
    dataloader: true
```

The properties and their relevant purposes are described below:

- `name`: This is only used to namespace the repository being created. It has no bearing on the generated code, as all repository generated this way assume that it is the only repository in the given package.
- `implements`: All repositories are expected to be based off any existing interface. Every method in this interface should use `context.Context` as its first parameter. The generated repository will implement this interface by creating the required for it. As such, this means that the `methods` array field should match every interface method name
- `methods`: Methods is an array of objects that describe the method that should be generated.
- `methods.[].name`: The name of the method. Every name should match one from the interface in `implements`.
- `methods.[].cypher`: The Cypher query that this method should use to lookup data.
- `methods.[].output`: The output that Cypher uses as part of its return query
- `methods.[].dataloader`: Determines whether a dataloader code generation strategy should be used.

It is expected that the interface methods provide named parameters:

```go
// Bad
type BadExample interface {
    Find(context.Context, string) (*model.Example, error)
}

// Good
type Repository interface {
    Find(ctx context.Context, id string) (*model.Example, error)
}
```

Additionally, the Cypher query expects that the second parameter be aligned with the parameters

```cypher
MATCH (grade:Grade) WHERE res.id = $id RETURN grade LIMIT 1
```

The code generator inspect the Cypher query and attempts to correlate its bindings (in this case `$id`) with the repository arguments. If this does not align, the generator will exit.

For example, taking the previous yaml example before:

```yaml
name: RoundRepository
implements: graphy/pkg/competition/rounds.Repository
methods:
  - name: FindRoundsByID
    cypher: |-
      MATCH (res:Round) WHERE res.id IN $roundIDs RETURN res LIMIT 10
    output: res
  - name: FindGradeRounds
    cypher: |-
      MATCH (g:Grade)-[:HAS_ROUND]->(r:Round)
      WHERE g.id IN $gradeIDs
      RETURN { id: g.id, items: COLLECT(r) } as out
    output: out
    dataloader: true
```

If the `FindRoundsByID()` interface method did _not_ have a `roundIDs []string` as part of its parameters, it would fail:

```go
// Assume the package path is graphy/pkg/competition/rounds
package rounds

import (
	"context"
	"graphy/transport/graphql/model"
)

type Repository interface {
	// This will fail autogeneration, the Cypher query does not take a `$id` binding or parameter!
	FindByRoundsID(ctx context.Context, id string) (*model.Round, error)
}
```