import { NgModule } from '@angular/core';
import { AuthComponent } from './auth.component';
import { AuthRoutingModule } from '../auth/auth-routing.module';
import { SharedModule } from '../shared';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

@NgModule({
	declarations: [
		AuthComponent
	],
	imports: [
		AuthRoutingModule,
		SharedModule,
		FormsModule,
		ReactiveFormsModule
	]
})
export class AuthModule { }
