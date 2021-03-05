package MixBooth;
use Dancer2;

our $VERSION = '0.1';
set serializer => 'JSON';

sub read_m3u {
	my ($path) = @_;

	open my $fh, "<", $path
		or return [];

	my @songs;
	while (<$fh>) {
		s/^\s+|\s+$//;
		next unless m/^(#\s*)?(.+)$/;

		push @songs, {
			active => $1 ? 0 : 1,
			file   => $2,
		};
	}
	close $fh;

	return \@songs;
}

sub write_m3u {
	my ($fh, $songs) = @_;

	for my $song (@$songs) {
		printf $fh "%s%s\n", $song->{active} ? '' : '#', $song->{file};
	}
}

get '/playlist' => sub {
	return { playlist => read_m3u("$ENV{RADIO_ROOT}/playlist.m3u") };
};

put '/playlist' => sub {
	my %in = params('body');

	open my $fh, ">", "$ENV{RADIO_ROOT}/.playlist.m3u"
		or return { error => 'unable to write playlist' };
	write_m3u($fh, $in{playlist});
	close $fh;
	rename "$ENV{RADIO_ROOT}/.playlist.m3u",
	       "$ENV{RADIO_ROOT}/playlist.m3u";

	return { ok => 'updated' };
};

post '/upload' => sub {
	my %in = params('body');

	my @cmd = ('ingest', $in{url});
	system(@cmd);
	return { ok => 'ingested' };
};

any '**' => sub {
	return { error => 'not found' };
};

true;
