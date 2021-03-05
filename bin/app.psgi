#!/usr/bin/env perl

use strict;
use warnings;
use Plack::Builder;
use FindBin;
use lib "$FindBin::Bin/../lib";
use MixBooth;
 
builder {
	enable "Plack::Middleware::Static",
	       root => './dist/',
	       path => sub { s|^/?$|/index.html|; m{[^/(playlist|upload)]} };
	MixBooth->to_app;
};
