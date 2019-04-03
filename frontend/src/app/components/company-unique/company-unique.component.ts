import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-company-unique',
  templateUrl: './company-unique.component.html',
  styleUrls: ['./company-unique.component.css']
})
export class CompanyUniqueComponent implements OnInit {
  code: String;

  constructor(private route: ActivatedRoute) {
     //code = this.route.snapshot.paramMap.get('code');
  }

  ngOnInit() {
    //this.code = this.route.snapshot.paramMap.get('code');
    this.code = this.route.snapshot.params.code;
  }

}
