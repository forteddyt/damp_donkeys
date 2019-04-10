import { Component, OnInit } from '@angular/core';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { Router, ActivatedRoute, NavigationExtras } from '@angular/router';

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
  code = ''

  constructor(private http: HttpClient, private router: Router, private route: ActivatedRoute) { }

  ngOnInit() {
    document.getElementById("employer_code").focus();
  }

  onKey(event: any)
  {
    this.code = event.target.value;
  }

  handleClick(event: any)
  {
    //this.code = "Clicked"
    this.router.navigateByUrl('/employers/' + this.code);
  }
}
