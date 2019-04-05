import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { Router, ActivatedRoute, NavigationStart } from '@angular/router';

@Component({
  selector: 'app-select-company',
  templateUrl: './select-company.component.html',
  styleUrls: ['./select-company.component.css']
})
export class SelectCompanyComponent implements OnInit {
  student: Observable<object>;

  constructor(public route: ActivatedRoute) {
    this.route.queryParams.subscribe(params => {
      console.log(params);
    });
  }

  ngOnInit() { 
  }
}
