import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
	selector: 'app-auth',
	templateUrl: './auth.component.html',
	styleUrls: ['./auth.component.css']
})
export class AuthComponent implements OnInit {
	authType: String = '';
	authForm: FormGroup;
	title: String = '';

	constructor(
		private route: ActivatedRoute,
		private fb: FormBuilder
	) {
		this.authForm = this.fb.group({
			'username': ['', Validators.required],
			'password': ['', Validators.required]
		});
	}

	ngOnInit(): void {
		this.route.url.subscribe(data => {
			this.authType = data[data.length - 1].path;
			this.title = (this.authType === 'login') ? 'Sign in' : 'Sign up';
		})
	}

	submitForm() {
		const credentials = this.authForm.value;
	}

}
