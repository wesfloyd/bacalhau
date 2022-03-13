
Prerequisites: 
- Create your Digital Ocean (DO) account, personal access token, passwordless SSH Key. Please see the prerequisites section [in this tutorial](https://www.digitalocean.com/community/tutorials/how-to-use-terraform-with-digitalocean#prerequisites) for further instructions
- [Install Terraform](https://www.terraform.io/downloads)

```
export DO_PAT="your_personal_access_token"

terraform init

#Confirm the plan works without errors
terraform plan \
  -var "do_token=${DO_PAT}" \
  -var "pvt_key=$HOME/.ssh/id_rsa" 

terraform apply \
  -var "do_token=${DO_PAT}" \
  -var "pvt_key=$HOME/.ssh/id_rsa"

```


Bring down terraform insfrastructure
```
terraform plan -destroy -out=terraform.tfplan \
  -var "do_token=${DO_PAT}" \
  -var "pvt_key=$HOME/.ssh/id_rsa"

terraform apply terraform.tfplan

```
