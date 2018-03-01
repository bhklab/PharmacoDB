FROM ruby:2.4

# see update.sh for why all "apt-get install"s have to stay as one long line
RUN apt-get update && apt-get install -y nodejs --no-install-recommends && rm -rf /var/lib/apt/lists/*

# see http://guides.rubyonrails.org/command_line.html#rails-dbconsole
RUN apt-get update && apt-get install -y mysql-client postgresql-client sqlite3 --no-install-recommends && rm -rf /var/lib/apt/lists/*

ENV RAILS_VERSION 5.1.3

RUN gem install rails --version "$RAILS_VERSION"

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# COPY Gemfile /usr/src/app/
# COPY Gemfile.lock /usr/src/app/
COPY . /usr/src/app/
RUN bundle install

COPY . /usr/src/app

EXPOSE 3000

RUN ["rake", "log:clear"]

RUN ["rm", "-rf", "tmp"]

# RUN ["bundle", "exec", "rake", "assets:precompile"]

# CMD ["bundle", "exec", "rails", "s", "-p", "3000", "-b", "'0.0.0.0'", "-e", "production"]
CMD ["rails", "server", "-b", "0.0.0.0", "-e", "production"]