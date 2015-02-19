PostGres
========

### Development

To run tests for the `postgres` subpackage enable the UUID extension:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

More info here: http://stackoverflow.com/a/12505220
