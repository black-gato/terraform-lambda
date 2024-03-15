resource "aws_iam_role" "lambda_role" {
  name               = "location_Lambda_function_Role"
  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "lambda.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF


}

resource "aws_iam_policy" "iam_policy_lambda" {
    name = "aws_iam_policy_for_terraform_aws_lambda_role"
    path = "/"
    description =  "AWS IAM Policy for managing aws lambda role"
  policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": [
       "logs:CreateLogGroup",
       "logs:CreateLogStream",
       "logs:PutLogEvents"
     ],
     "Resource": "arn:aws:logs:*:*:*",
     "Effect": "Allow"
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role" {
    role = aws_iam_role.lambda_role.name
    policy_arn = aws_iam_policy.iam_policy_lambda.arn
  
}




resource "aws_lambda_function" "hello-world" {
    function_name = "hello-world"
    filename = "hello.zip"
    handler = "hello"
    description = "playing around"
    role = aws_iam_role.lambda_role.arn
    source_code_hash = filebase64sha256("hello.zip")
    
    memory_size = 128
    runtime = "provided.al2023"

    
    depends_on = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]

  
}