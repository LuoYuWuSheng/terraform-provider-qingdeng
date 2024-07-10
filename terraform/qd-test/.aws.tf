provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1f0" # 这是一个 Amazon Linux 2 AMI 的示例 ID，请根据需要更改
  instance_type = "t2.micro"

  tags = {
    Name = "example-instance"
  }
}
