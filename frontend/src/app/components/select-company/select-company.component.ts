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

  constructor(private router: Router) {
    console.log(this.router.getCurrentNavigation().extras.state);
  }

  ngOnInit() { 
  }
}
