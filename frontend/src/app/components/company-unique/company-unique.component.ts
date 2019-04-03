import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-company-unique',
  templateUrl: './company-unique.component.html',
  styleUrls: ['./company-unique.component.css']
})
export class CompanyUniqueComponent implements OnInit {
  code: String;
  validString: String = "abcdefghij";

  constructor(private route: ActivatedRoute, private router: Router) {
     //code = this.route.snapshot.paramMap.get('code');
  }

  ngOnInit() {
    //this.code = this.route.snapshot.paramMap.get('code');
    this.code = this.route.snapshot.params.code;
    if (this.code === this.validString) {}
    else {
      this.router.navigate([''], {});
    }
  }

}
