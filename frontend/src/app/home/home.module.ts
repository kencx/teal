import { NgModule } from '@angular/core';
import { FormsModule } from "@angular/forms";
import { CommonModule } from '@angular/common';
import { HomeComponent } from './home.component';
import { BooksComponent } from '../books/books.component';
import { SharedModule } from "../shared";

@NgModule({
  declarations: [
    HomeComponent,
    BooksComponent
  ],
  imports: [
    CommonModule,
	SharedModule,
	FormsModule
  ]
})
export class HomeModule { }
