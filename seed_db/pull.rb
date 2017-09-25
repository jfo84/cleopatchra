require_relative './data_mixin'

class Pull
  include Virtus.model
  include DataMixin

  attribute :repo_id, String

  def record
    connection.exec("INSERT INTO pulls (id, data, repo_id) VALUES ($1, $2, $3)", [id, data, repo_id])
  end

  def url
    data_hash['url']
  end

  def is_dup?
    result = connection.exec("SELECT id FROM pulls WHERE id = $1", [id])
    result.ntuples == 1
  end
end