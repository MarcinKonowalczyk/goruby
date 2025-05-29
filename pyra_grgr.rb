# TRANSFORM_PIPELINE *ast.Assignment
__lifted_function_with_a_comment_in_it = (-> () {# i'm a little comment;return "This function has a comment in it."})

# TRANSFORM_PIPELINE *ast.FunctionLiteral
-> () {# i'm a little comment;return "This function has a comment in it."}

# TRANSFORM_PIPELINE *ast.ContextCallExpression
puts(__lifted_function_with_a_comment_in_it[])

# => "This function has a comment in it."
