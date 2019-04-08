import { Component, OnInit } from '@angular/core';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';

@NgModule({
  imports: [
    CommonModule,
    BrowserModule,
    FormsModule
  ],
  declarations: []
})

@Component({
  selector: 'app-company-welcome',
  templateUrl: './company-welcome.component.html',
  styleUrls: ['./company-welcome.component.css']
})
export class CompanyWelcomeComponent implements OnInit {
  code: String

  constructor(private route: ActivatedRoute, private router: Router) { }

  ngOnInit() {
    //this.code = this.route.snapshot.params.input;
    //if (this.code === "")
    //{
    //  this.router.navigate(['employers/something'], {});
    //}
  }

}
