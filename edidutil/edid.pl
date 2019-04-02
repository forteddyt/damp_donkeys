#!/usr/bin/perl

use strict;
use IO::Socket::SSL;
use Authen::SASL;
use Net::LDAP;

use Data::Dumper;

my %cfg;
$cfg{'ED_ID_HOST'}	= 'id.directory.vt.edu';
$cfg{'ED_ID_PORT'}      = 389;
$cfg{'ED_ID_VER'}       = 3;
$cfg{'ED_ID_BASE'}      = 'ou=people,dc=vt,dc=edu';
$cfg{'ED_ID_CA'} 	= 'ed_id_ca.pem';
$cfg{'ED_ID_CERT'} 	= 'public.crt';
$cfg{'ED_ID_KEY'} 	= 'private.key';
$cfg{'ED_ID_USER'} 	= 'dn:uusid=cs-career-fair,ou=services,dc=vt,dc=edu';

sub ed_id
{
    my %args = @_;
    my $filter_val = $args{'FILTER_VAL'} or return( undef() );
    my $filter_key = $args{'FILTER_KEY'} || 'virginiaTechID';
    my $attribref = $args{'OUTPUT'} || undef();

    my $base = $args{'LDAP_BASE'} || $cfg{'ED_ID_BASE'};
    my $host = $args{'LDAP_HOST'} || $cfg{'ED_ID_HOST'};
    my $port = $args{'LDAP_PORT'} || $cfg{'ED_ID_PORT'};
    my $user = $args{'LDAP_USER'} || $cfg{'ED_ID_USER'};
    
    my $sasl = Authen::SASL->new(mechanism => 'EXTERNAL', callback => {pass => '', user => $user});

    # retry a few times if the first connection failed
    my $ldap = Net::LDAP->new( $host, port => $port, version => $cfg{'ED_ID_VER'} );
    if ( ! $ldap )
    {
	sleep( 1 );
	$ldap = Net::LDAP->new( $host, port => $port, version => $cfg{'ED_ID_VER'} );
	if ( ! $ldap )
	{
	  sleep ( 1 );
	  $ldap = Net::LDAP->new( $host, port => $port, version => $cfg{'ED_ID_VER'} );
	  if ( ! $ldap )
	    { return( undef() ); }
	}
    }

    my $mesg = $ldap->start_tls(verify => 'require', cafile => $cfg{'ED_ID_CA'}, clientcert => $cfg{'ED_ID_CERT'},
                            clientkey => $cfg{'ED_ID_KEY'});
    if ( $mesg->code )
    {
      $ldap->unbind;
      return( undef() );
    }
    $mesg = $ldap->bind(dn => '', sasl => $sasl);
    if ( $mesg->code )
    {
      $ldap->unbind;
      return( undef() );
    }

    if ( $attribref )
	{ $mesg = $ldap->search( base => $base, filter => "(&(!(eduPersonAffiliation=VT-GUEST))($filter_key=$filter_val))", attrs => $attribref ); }
    else
	{ $mesg = $ldap->search( base => $base, filter => "(&(!(eduPersonAffiliation=VT-GUEST))($filter_key=$filter_val))" ); }
    $ldap->unbind;

    return( undef() ) if ( $mesg->code );
    return( $mesg->as_struct );    
}

my $id = shift;
my $ref = ed_id(FILTER_VAL => $id);
my ($key, $hashRef) = each(%$ref); 
if ( $hashRef )
{
#  print Dumper($hashRef);
   print( join(';', $hashRef->{displayname}[0], $hashRef->{major}[0], $hashRef->{classlevel}[0]) );
   exit 0;
}
exit 1;

