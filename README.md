Examples
========
Here are a few examples to get you started!

Fetch all users
```go
query := gosql.Select().
    From("users", []string{"id", "first_name"})
```

Results in
```sql
SELECT id, first_name FROM users
```
________

Fetch users and their payments who are active and have amounts between $10 and $20
```go
query := gosql.Select().
    From("users", []string{"id", "first_name").
    InnerJoin("payments", "payments.user_id = users.id", []string{"amount", "is_approved"}).
    Where("payments.amount BETWEEN ? AND ? AND users.is_active", 10, 20)
```

Results in
```sql
SELECT id, first_name, amount, is_approved FROM users 
    INNER JOIN payments ON payments.user_id = users.id 
    WHERE payments.amount BETWEEN ? AND ? AND users.is_active
```