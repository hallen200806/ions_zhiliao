package test_tree



[
	{
		"auth_name":"用户管理",
		"url_for":"#",
		"children":[
			{
				"auth_name":"用户列表",
				"url_for":"UserController.List",
				"children":[]
			},
			{
				"auth_name":"用户列表1",
				"url_for":"UserController.List1",
				"children":[]
			}

		]
	},
	{
		"auth_name":"权限管理",
		"url_for":"#",
		"children":[
			{
				"auth_name":"权限列表",
				"url_for":"AuthController.List",
			},
			{
				"auth_name":"权限列表1",
				"url_for":"AuthController.List1"
			}

		]

	}

]
