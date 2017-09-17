# Terraform SDK

In an effort to vendor the `Terraform` code base directly into another Go program I have written a (somewhat hacky) SDK for `Terraform`.

It is important to note that this is an unofficial port of the program and absolutely nothing is supported.

A lot of the code in the `tfmain` directory is **copypasta** from the terraform code base, and was pulled into a new package for ancillary reasons in getting the code to vendor smoothly.
It is important to note that this a **high level abstraction** of the `Terraform` code base. 
The SDK supports hooking into the program in the same way a user would interact with the program via the CLI, except through idiomatic Go.


The original code base (and it's fabulous Mozilla Public License 2.0) can be found [here](https://github.com/hashicorp/terraform).
A huge thanks to the amazing people at Hashicorp and the open source Terraform community for giving us this amazing open source tool. 
We are all very grateful!