import { Component, OnInit } from '@angular/core';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: []
})

@Component({
  selector: 'app-company-welcome',
  templateUrl: './company-welcome.component.html',
  styleUrls: ['./company-welcome.component.css']
})
export class CompanyWelcomeComponent implements OnInit {

  constructor() { }

  ngOnInit() {
  }

}
