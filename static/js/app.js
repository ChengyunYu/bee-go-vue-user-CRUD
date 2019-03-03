(function(Vue) {
	"use strict";

	const recordPerPage = 10;

	const FAILURE = "FAILURE";
	const SUCCESS = "SUCCESS";

	let validateEmail = (rule, value, callback) => {
		let re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
		if (!re.test(value)) {
			callback(new Error("Please enter a valid email address"))
		} else {
			callback()
		}
	};

	let validatePassword1 = (rule, value, callback) => {
		if (value === '') {
			callback(new Error("Please enter a password. "));
		} else {
			if (vue_instance.newUser.password_2 !== '') {
				vue_instance.$refs.newUser.validateField('password_2');
			}
			callback();
		}
	};

	let validatePassword2 =  (rule, value, callback) => {
		if (value === '') {
			callback(new Error("Please re-enter password. "));
		} else if (value !== vue_instance.newUser.password_1) {
			callback(new Error("The two passwords does not match. "));
		} else {
			callback();
		}
	};

    let validateNewPassword1 = (rule, value, callback) => {
        if (value === '') {
            callback(new Error("Please enter a password. "));
        } else {
            if (vue_instance.newPassword.password_2 !== '') {
                vue_instance.$refs.newPassword.validateField('password_2');
            }
            callback();
        }
    };

    let validateNewPassword2 =  (rule, value, callback) => {
        if (value === '') {
            callback(new Error("Please re-enter password. "));
        } else if (value !== vue_instance.newPassword.password_1) {
            callback(new Error("The two passwords does not match. "));
        } else {
            callback();
        }
    };

	var vue_instance = new Vue({
		el: '#app',

		data: {
			keyword: "",
            totalPage: 1,
            currentPage: 1,
            totalRecords: 0,
			searchUsers: [],
			currentUser: {
                surname: "",
                given_name: "",
                email: "",
                org: "",
			},
			editedUser: {
				surname: "",
				given_name: "",
				email: "",
				org: "",
			},
            newPassword: {
                password_1: "",
                password_2: "",
            },
			newUser: {
				surname: "",
				given_name: "",
				password_1: "",
				password_2: "",
				email: "",
				org: "",
			},
			create_rules: {
				surname: [
					{ required: true, message: 'Please enter surname. ', trigger: 'blur' },
				],
				given_name: [
					{ required: true, message: 'Please enter given name. ', trigger: 'blur' },
				],
				password_1: [
					{ validator: validatePassword1, trigger: 'blur'},
				],
				password_2: [
					{ validator: validatePassword2, trigger: 'blur'},
				],
				email: [
					{ required: true, message: 'Please enter email address. ', trigger: 'blur' },
					{ validator: validateEmail, trigger: 'blur' },
				],
				org: [
					{ required: true, message: 'Please select an organization. ', trigger: 'change' },
				],
			},
            edit_rules: {
                surname: [
                    { required: true, message: 'Please enter surname. ', trigger: 'blur' },
                ],
                given_name: [
                    { required: true, message: 'Please enter given name. ', trigger: 'blur' },
                ],
                email: [
                    { required: true, message: 'Please enter email address. ', trigger: 'blur' },
                    { validator: validateEmail, trigger: 'blur' },
                ],
                org: [
                    { required: true, message: 'Please select an organization. ', trigger: 'change' },
                ],
            },
            password_rules: {
                password_1: [
                    { validator: validateNewPassword1, trigger: 'blur'},
                ],
                password_2: [
                    { validator: validateNewPassword2, trigger: 'blur'},
                ],
            },
			orgs: [
				{ text: "Engineers", value: "engineers" },
				{ text: "Managers", value: "managers" },
			],
            dialogEditVisible: false,
			dialogVisible: false,
            dialogPasswordVisible: false,
			loadingCreate: false,
			loadingEdit: false,
			loadingTable: false,
			loadingPassword: false,
			hackReset: true,
		},

		created: function () {
		},

		// Functions we will be using
		methods: {
			clearCreateForm: function () {
				this.getUsers();
				this.$http.get('/users/recent').then(function (res) {
					this.recentUsers = res.data.iterms ? res.data.iterms : [];
				});
			},

			closeCreateDialog() {
				this.dialogVisible=false;
				this.$refs.newUser.resetFields();
			},

            closeEditDialog() {
                this.dialogEditVisible=false;
                this.$refs.editedUser.resetFields();
            },

            closePasswordDialog() {
                this.dialogPasswordVisible=false;
                this.$refs.newPassword.resetFields();
            },

			closeForm(formName) {
			    if(formName === "newUser") {
                    this.dialogVisible = false;
                }
			    else if(formName === "currentUser") {
			        this.dialogEditVisible = false;
                }
			    else {
                    this.dialogPasswordVisible = false;
                }
				this.$refs[formName].resetFields();
			},

			resetForm(formName) {
				this.$refs[formName].resetFields();
			},

			clearInput: function () {
				this.keyword = "";
			},

			getUsers: function () {
				let keywordPart = "";
				if(this.keyword !== "") {
					keywordPart = this.keyword + "/"
				}
				this.$http.get('/users/' + keywordPart + this.currentPage).then(function (res) {
					console.log(res);
					if(res.body.result === SUCCESS) {
						this.searchUsers = res.body.results.users;
						this.totalPage = res.body.total_page;
						this.totalRecords = res.body.total_results;
						this.currentPage = res.body.current_page;
					} else {
						this.$notify.error({
							title: "Failed to retrieve latest user(s): ",
							message: res.body.message,
						});
					}
					this.loadingTable = false;
				})
			},

			getUsersFirst: function () {
				this.currentPage = 1;
				this.getUsers();
			},

			getUsersPage: function () {
				this.getUsers();
			},

            createUser: function () {
                this.$refs.newUser.validate((valid) => {
                    if (valid) {
                    	this.loadingCreate = true;
                        this.$http.post('/user',this.newUser).then(function(res) {
                            if(res.body.result === FAILURE) {
								this.$notify.error({
									title: "Failed to add user: ",
									message: res.body.message,
								});
								this.loadingCreate = false;
							}
                            else {
								this.$notify({
									title: "Successfully added user: ",
									message: res.body.message,
									type: 'success',
								});
								this.$refs.newUser.resetFields();
								this.loadingCreate = false;
								this.dialogVisible = false;
								this.getUsersPage()
							}
                        });
                    } else {
                        return false;
                    }
                });
            },

			showCreateUser: function () {
				this.dialogVisible = true;
			},

            showEditUser: function(row) {
				this.hackReset = false;
				this.$nextTick(() => {
					this.hackReset = true;
					this.editedUser.surname = row.surname;
					this.editedUser.given_name = row.given_name;
					this.editedUser.email = row.email;
					this.editedUser.org = row.org;
					this.currentUser.surname = row.surname;
					this.currentUser.given_name = row.given_name;
					this.currentUser.email = row.email;
					this.currentUser.org = row.org;
					this.dialogEditVisible = true;
				});
            },

			editUser: function() {
				this.$refs.editedUser.validate((valid) => {
					if (valid) {
						this.loadingEdit = true;
						this.$http.put('/user/' + this.currentUser.email + "/" + this.currentUser.org + "/edit",
							this.editedUser).then(function(res) {
							if(res.body.result === FAILURE) {
								this.$notify.error({
									title: "Failed to update user: ",
									message: res.body.message,
								});
								this.loadingCreate = false;
							}
							else {
								this.$notify({
									title: "Successfully updated user: ",
									message: res.body.message,
									type: 'success',
								});
								this.$refs.editedUser.resetFields();
								this.loadingEdit = false;
								this.dialogEditVisible = false;
								this.getUsersPage();
							}
						});
					} else {
						return false
					}
				});
			},

			updatePassword: function() {
				this.$refs.newPassword.validate((valid) => {
					if (valid) {
						this.loadingPassword = true;
						this.$http.put('/user/' + this.currentUser.email + "/" + this.currentUser.org + "/password",
							this.newPassword).then(function (res) {
							if(res.body.result === SUCCESS) {
								this.$notify({
									title: "Successfully updated password for user: ",
									message: res.body.message,
									type: 'success',
								});
								this.dialogPasswordVisible = false;
							} else {
								this.$notify.error({
									title: "Failed to updated password for user: ",
									message: res.body.message,
								});
							}
							this.loadingPassword = false;
						});
					}
				});
			},

            showUpdatePassword: function() {
			    this.dialogPasswordVisible = true;
            },

			deleteUser: function (row) {
				this.$confirm('Are you sure you want to delete ' + row.email + "? ")
					.then(_ => {
						this.$http.delete('/user/' + row.email + "/" + row.org).then(function (res) {
							if(res.body.result === SUCCESS) {
								this.$notify({
									title: "Successfully deleted user: ",
									message: res.body.message,
									type: 'success',
								});

								this.getUsersPage();
							}
							else {
								this.$notify.error({
									title: "Failed to delete user: ",
									message: res.body.message,
								});
							}
						})
					})
					.catch(_ => {});
			},
		}
	});
})(Vue);
