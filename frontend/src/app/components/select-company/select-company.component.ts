import { Component, AfterViewInit, ComponentFactoryResolver, ViewChild, ViewContainerRef } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router} from '@angular/router';
import { SelectCompanyTileComponent } from '../select-company-tile/select-company-tile.component';


@Component({
  selector: 'app-select-company',
  templateUrl: './select-company.component.html',
  styleUrls: ['./select-company.component.css']
})

export class SelectCompanyComponent implements AfterViewInit {
  @ViewChild('companyInsert', { read: ViewContainerRef }) companyInsert: ViewContainerRef;
  stateData
  companyList
  header
  sticky

  constructor(private http: HttpClient, private router: Router, private vcr: ViewContainerRef, private cfr: ComponentFactoryResolver) {
    this.stateData = this.router.getCurrentNavigation().extras.state;
    if(this.stateData == null || this.stateData.pid == null){
      console.log("Invalid state has been passed (or not passed at all); should redirect to '/'");
      this.router.navigate(['/']);
    }
  }

  async ngAfterViewInit() {
    const resp = await this.getCompanies();
    
    console.log(resp.status)
    this.companyList = resp.body
    this.loadComponents()
  }

  getCompanies() {
    return this.http.get("https://csrcint.cs.vt.edu/api/company_list", {observe: 'response'}).toPromise();
  }

  loadComponents() {
    console.log(this.companyList)
    const cFactory = this.cfr.resolveComponentFactory(SelectCompanyTileComponent);
    this.companyInsert.clear();
    for(var i in this.companyList){
      const companyComponent = <SelectCompanyTileComponent>this.companyInsert.createComponent(cFactory).instance;
      
      companyComponent.companyName = this.companyList[i];
      companyComponent.stateData = this.stateData;
    }
  }
}
