import { Component, AfterViewInit, ComponentFactoryResolver, ViewChild, ViewContainerRef } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router, ActivatedRoute, NavigationStart } from '@angular/router';
import { SelectCompanyTileComponent } from "../select-company-tile/select-company-tile.component"

@Component({
  selector: 'app-select-company',
  templateUrl: './select-company.component.html',
  styleUrls: ['./select-company.component.css']
})
export class SelectCompanyComponent implements AfterViewInit {
  @ViewChild('companyInsert', { read: ViewContainerRef }) companyInsert: ViewContainerRef;
  state_vars
  company_list
  header
  sticky

  constructor(private http: HttpClient, private router: Router, private vcr: ViewContainerRef, private cfr: ComponentFactoryResolver) {
    this.state_vars = this.router.getCurrentNavigation().extras.state;
    if(this.state_vars == null || this.state_vars.pid == null){
      console.log("null state; should redirect to '/'")
    }
  }

  async ngAfterViewInit() {
    const resp = await this.getCompanies();
    
    console.log(resp.status)
    this.company_list = resp.body
    this.loadComponents()
  }

  getCompanies() {
    return this.http.get("https://csrcint.cs.vt.edu/api/company_list", {observe: 'response'}).toPromise();
  }

  loadComponents() {
    console.log(this.company_list)
    const cFactory = this.cfr.resolveComponentFactory(SelectCompanyTileComponent);
    this.companyInsert.clear();
    for(var i in this.company_list){
      const companyComponent = <SelectCompanyTileComponent>this.companyInsert.createComponent(cFactory).instance;
      companyComponent.companyName = this.company_list[i];
    }
  }
}
