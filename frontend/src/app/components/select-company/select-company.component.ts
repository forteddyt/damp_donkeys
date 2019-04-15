import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { Router, ActivatedRoute, NavigationExtras } from '@angular/router';

@Component({
  selector: 'app-select-company',
  templateUrl: './select-company.component.html',
  styleUrls: ['./select-company.component.css']
})
export class SelectCompanyComponent implements OnInit {
  student: Observable<object>;
  companies: Object;
  constructor(public route: Router, private http: HttpClient) {
    var state = this.route.getCurrentNavigation().extras.state;
    console.log(state);
    this.http.get("https://csrcint.cs.vt.edu/api/company_list").subscribe((res) => {
      console.log(res);
      this.companies = res;
      console.log(this.companies)
    });
  }

  ngOnInit() { 
  }
}
